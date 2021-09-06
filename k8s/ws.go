package k8s

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// 升级http请求为websocket协议
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
