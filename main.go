package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"regexp"
	"time"
)

var (
	lackOfResources string = ".*Insufficient.*cpu.*Insufficient.*memory.*"
	mountFailed     string = ".*Unable.*to.*attach.*or.*mount.*volumes.*unattached.*"
	startupFailed   string = ".*runc.*create.*failed.*unable.*start.*container.*process.*exec.*not.*found.*in.*unknown.*"
)

func main() {
	kubeconfig := "/Users/huangjianchen/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, _ := kubernetes.NewForConfig(config)
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*3)
	evInformer := informerFactory.Core().V1().Events()
	//informer := podInformer.Informer()

	evLister := evInformer.Lister()

	//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    onAdd,
	//	UpdateFunc: onUpdate,
	//	DeleteFunc: onDelete,
	//})

	stopCh := make(chan struct{})
	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	defer close(stopCh)

	event, _ := evLister.Events("").List(labels.Everything())
	for _, events := range event {
		// 启动失败
		if events.Count >= 5 {
			if events.Reason == "Failed" {
				if matched, _ := regexp.Match(startupFailed, []byte(events.Message)); matched {
					log.Printf("发现异常! 资源类型:%s 异常种类:启动失败 应用名字:%s 集群:%s 原因分析:可能存在启动命令配置错误", events.InvolvedObject.Kind, events.Name, events.Namespace)
				}
			} else if events.Reason == "FailedScheduling" {
				if matched, _ := regexp.Match(lackOfResources, []byte(events.Message)); matched {
					log.Printf("发现异常! 资源类型:%s 异常种类:调度失败 应用名字:%s 集群:%s 原因分析:集群CPU或内存资源或者内存资源不够 ", events.InvolvedObject.Kind, events.Name, events.Namespace)
				}
			} else if events.Reason == "FailedMount" {
				if matched, _ := regexp.Match(mountFailed, []byte(events.Message)); matched {
					log.Printf("发现异常! 资源类型:%s 异常种类:启动失败 应用名字:%s 集群:%s 原因分析:存储卷挂载失败或无法挂载 ", events.InvolvedObject.Kind, events.Name, events.Namespace)
				}
			}
		}
	}
}

func onAdd(obj interface{}) {
	deploy := obj.(*v1.Pod)
	fmt.Println("add a deployment:", deploy.Name)
	//fmt.Println(obj)
}

func onUpdate(old, new interface{}) {
	oldDeploy := old.(*v1.Pod)
	newDeploy := new.(*v1.Pod)
	fmt.Println("update deployment:", oldDeploy.Name, newDeploy.Name)
	//fmt.Println(new)
}

func onDelete(obj interface{}) {
	deploy := obj.(*v1.Pod)
	fmt.Println("delete a deployment:", deploy.Name)
}
