package main

import (
	"eventWatchController/we"
	"fmt"
	"os"
)

func main() {
	kubeconfig := fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.kube/config")
	eventWatch := we.NewInformer(kubeconfig)
	events := eventWatch.CreateEventInformer()
	eventWatch.Run(events, "Back-off restarting failed container.*", "启动失败", "容器重启失败，命令执行错误请检查")
	eventWatch.Run(events, ".*Insufficient.*cpu.*Insufficient.*memory.*", "启动失败", "容器需要资源过多，节点无法满足 请检查配额")
	eventWatch.Run(events, ".*Unable.*to.*attach.*or.*mount.*volumes.*unattached.*", "无法调度", "存储卷不存在，请创建后重试")
	eventWatch.Run(events, ".*runc.*create.*failed.*unable.*start.*container.*process.*exec.*not.*found.*in.*unknown.*", "启动失败", "容器命令不存在，请确认命令是否存在于容器之中")

	<-we.We.StopCh
	defer close(we.We.StopCh)
}
