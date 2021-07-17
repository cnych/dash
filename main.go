package main

import (
	"flag"

	"github.com/cnych/dash/global"
	"github.com/cnych/dash/routers"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func main() {
	// todo，传递解析flag参数
	// 初始化 klog，也可以绑定到本地的flagset
	klog.InitFlags(nil)
	defer klog.Flush()
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()

	// 全局初始化
	if err := global.Init(); err != nil {
		klog.V(2).ErrorS(err, "init global failed")
		return
	}
	serv := gin.Default()
	// 注册路由
	routers.InitApi(serv)
	// 启动服务 0.0.0.0:8888
	if err := serv.Run(":8888"); err != nil {
		klog.V(2).ErrorS(err, "server run failed")
	}
}
