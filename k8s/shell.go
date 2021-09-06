package k8s

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type ShellMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
}

type KubeShell struct {
	conn *websocket.Conn
	sizeChan   chan remotecommand.TerminalSize
	stopChan chan struct{}
	tty bool
}

var EOT = "\u0004"
var _ TtyHandler = &KubeShell{}

func NewKubeShell(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*KubeShell, error) {
	// 升级get请求为websocket协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	kubeShell := &KubeShell{
		conn: conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		stopChan: make(chan struct{}),
		tty: true,
	}
	return kubeShell, nil
}

func (k *KubeShell) Stdin() io.Reader {
	return k
}

func (k *KubeShell) Stdout() io.Writer {
	return k
}

func (k *KubeShell) Stderr() io.Writer {
	return k
}

func (k *KubeShell) Tty() bool {
	return k.tty
}

func (k *KubeShell) Next() *remotecommand.TerminalSize {
	select {
	case size := <- k.sizeChan:
		return &size
	case <-k.stopChan:
		return nil
	}
}

func (k *KubeShell) Done() {
	close(k.stopChan)
}

func (k *KubeShell) Close() error {
	return k.conn.Close()
}

func (k *KubeShell) Read(p []byte) (n int, err error) {
	_, message, err := k.conn.ReadMessage()
	if err != nil {
		return copy(p, EOT), err
	}
	var msg ShellMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return copy(p, EOT), err
	}
	switch msg.Type {
	case "read":
		return copy(p, msg.Data), nil
	case "resize":
		k.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return copy(p, EOT), fmt.Errorf("unknown message type: %s", msg.Type)
	}
}


func (k *KubeShell) Write(p []byte) (n int, err error) {
	msg, err := json.Marshal(ShellMessage{
		Type: "write",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}
	if err := k.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

