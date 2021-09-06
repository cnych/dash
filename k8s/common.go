package k8s

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// kubeclient 定义包含所有需要操作的client
type kubeclient struct {
	Pod  *PodClient
	Node *NodeClient
}

var Client *kubeclient

func initK8sClient() (*kubernetes.Clientset, *rest.Config, error) {
	var err error
	var config *rest.Config
	var clientset *kubernetes.Clientset
	// inCluster（Pod）、Kubeconfig（kubectl）
	// 通过flag传递kubeconfig参数
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "ydzs-config")
	// 首先使用 inCluster 模式（RBAC -> list|get node）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 kubeconfig 模式
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, nil, err
		}
	}

	// 创建clientset对象
	if clientset, err = kubernetes.NewForConfig(config); err != nil {
		return nil, nil, err
	}
	return clientset, config, nil
}

func NewKubeClient() error {
	clientset, config, err := initK8sClient()
	if err != nil {
		return err
	}
	Client = &kubeclient{
		Pod:  NewPodClient(clientset, config),
		Node: NewNodeClient(clientset),
	}
	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}