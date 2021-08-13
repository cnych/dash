package k8s

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type KubeLogger struct {
	Conn *websocket.Conn
}

func NewKubeLogger(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*KubeLogger, error) {
	// 升级get请求为websocket协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	kubeLogger := &KubeLogger{
		Conn: conn,
	}
	return kubeLogger, nil
}

func (kl *KubeLogger) Write(data []byte) error {
	return kl.Conn.WriteMessage(websocket.TextMessage, data)
}

func (kl *KubeLogger) Close() error {
	return kl.Conn.Close()
}
