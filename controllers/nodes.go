package controllers

import (
	"github.com/cnych/dash/k8s"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func GetNodeList(c *gin.Context) {
	nodes, err := k8s.Client.Node.List("")
	if err != nil {
		klog.V(2).ErrorS(err, "get node list failed", "controller", "GetNodeList")
		writeError(c, err.Error())
		return
	}
	writeOK(c, gin.H{"nodes": nodes})
}
