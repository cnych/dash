package k8s

import (
	"bufio"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type PodClient struct {
	clientset *kubernetes.Clientset
}

func NewPodClient(clientset *kubernetes.Clientset) *PodClient {
	return &PodClient{
		clientset: clientset,
	}
}

func (cli *PodClient) Get(name, namespace string) (*corev1.Pod, error)  {
	opt := metav1.GetOptions{}
	return cli.clientset.CoreV1().Pods(namespace).Get(context.Background(), name, opt)
}

func (cli *PodClient) Logs(name, namespace string, opts *corev1.PodLogOptions) *restclient.Request  {
	return cli.clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
}

func (cli *PodClient) LogsStream(name, namespace string, opts *corev1.PodLogOptions, writer io.Writer) error  {
	req := cli.Logs(name, namespace, opts)
	stream ,err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer stream.Close()

	buf := bufio.NewReader(stream)
	for {  // 一直从buffer中读取数据去
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_, err = writer.Write(bytes)
			}
			return err
		}
		_, err = writer.Write(bytes)
		if err != nil {
			return err
		}
	}

}