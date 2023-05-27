package we

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

//func NewInformer(kubeconfig *rest.Config, clientset *kubernetes.Clientset, indexinformer *cache.SharedIndexInformer) *weController {

func NewInformer(kubeconfig string) *weController {
	return &weController{
		Kubeconfig: kubeconfig,
	}
}

func (c *weController) CreateEventInformer() []*v1.Event {
	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, _ := kubernetes.NewForConfig(config)
	// 注册informer工厂 持续watchAPIserver 缓存同步时间 3分钟
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*3)
	// 获取kubernetes event资源类型
	evInformer := informerFactory.Core().V1().Events()
	// list API可以查询当前的资源及其对应的状态(即期望的状态)
	evLister := evInformer.Lister()

	// 启动informer工厂
	informerFactory.Start(c.StopCh)
	// 等待缓存同步完成
	informerFactory.WaitForCacheSync(c.StopCh)
	event, _ := evLister.Events("").List(labels.Everything())

	return event
}
