package main

import (
	"eventWatchController/we"
	"fmt"
	"os"
)

func main() {
	kubeconfig := fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.kube/config")
	eventWatch := we.NewInformer(kubeconfig)
	eventPod := eventWatch.CreateEventInformer()
	eventPod.Run("Back-off restarting failed container.*", "重启失败", "重启失败，命令执行错误请检查")
	eventPod.Run(".*Insufficient.*cpu.*Insufficient.*memory.*", "资源不足", "容器需要资源过多，节点无法满足 请检查配额")
	eventPod.Run(".*Unable.*to.*attach.*or.*mount.*volumes.*unattached.*", "存储卷", "存储卷不存在，请创建后重试")
	eventPod.Run(".*runc.*create.*failed.*unable.*start.*container.*process.*exec.*not.*found.*in.*unknown.*", "创建失败", "容器命令不存在，请确认命令是否存在于容器之中")
	<-we.We.StopCh
	defer close(we.We.StopCh)
}
