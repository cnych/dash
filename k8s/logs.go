package k8s

import (
	"net/http"

	"github.com/gorilla/websocket"
)


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

//Write(p []byte) (n int, err error)
func (kl *KubeLogger) Write(data []byte) (int, error) {
	if err := kl.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (kl *KubeLogger) Close() error {
	return kl.Conn.Close()
}
