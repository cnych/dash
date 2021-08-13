package controllers

import (
	"net/http"
	"strconv"

	"github.com/cnych/dash/k8s"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// GetKubeLogs 实时获取Pod的日志
func GetKubeLogs(c *gin.Context) {
	// /api/v1/namespaces/kube-system/pods/traefik-7cb4cb6bf5-p779x/logs?tailLines=500&timestamps=true&previous=false&container=traefik'
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.Query("container")
	tailLines, _ := strconv.ParseInt(c.DefaultQuery("tailLines", "500"), 10, 64)
	timestamps, _ := strconv.ParseBool(c.DefaultQuery("timestamps", "true"))
	previous, _ := strconv.ParseBool(c.DefaultQuery("previous", "false"))

	klog.V(2).InfoS("get kube logs request params", "namespace", namespace, "pod", podName,
		"container", container, "tailLines", tailLines, "timestamps", timestamps, "previous", previous)

	if namespace == "" || podName == "" || container == "" {
		c.String(http.StatusBadRequest, "must specific namespace、pod and container query params")
		return
	}

	// 获取pod的日志（websocket）
	// 把当前的get http request -> upgrade websocket
	kubeLogger, err := k8s.NewKubeLogger(c.Writer, c.Request, nil)
	if err != nil {
		klog.Error(err, "upgrade websocket failed")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 构造获取日志的结构体
	opts := corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		TailLines:  &tailLines,
		Timestamps: timestamps,
		Previous:   previous,
	}
	if err := k8s.Client.Pod.LogsStream(podName, namespace, &opts, kubeLogger); err != nil {
		klog.Error(err, "GetLogs stream failed")
		_, _ = kubeLogger.Write([]byte(err.Error()))
	}
}
