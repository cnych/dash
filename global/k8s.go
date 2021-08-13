package global

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	client *kubernetes.Clientset
)

func initK8sClient() error {
	var err error
	var config *rest.Config
	// inCluster（Pod）、Kubeconfig（kubectl）
	// 通过flag传递kubeconfig参数
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "ydzs-config")
	// 首先使用 inCluster 模式（RBAC -> list|get node）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 kubeconfig 模式
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return err
		}
	}

	// 创建clientset对象
	if client, err = kubernetes.NewForConfig(config); err != nil {
		return err
	}
	return nil
}

func K8sClient() *kubernetes.Clientset {
	return client
}
