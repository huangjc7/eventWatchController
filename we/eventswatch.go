package we

import (
	v1 "k8s.io/api/core/v1"
	"log"
	"regexp"
)

func (c *weController) Run(events []*v1.Event, regexFailed string, failedType string, reason string) {
	for _, event := range events {
		if event.Count >= 3 {
			// 启动失败
			if event.Reason == "BackOff" || event.Reason == "Failed" {
				if matched, _ := regexp.Match(regexFailed, []byte(event.Message)); matched {
					log.Printf("发现异常! 资源类型:%s 异常种类:%s 应用名字:%s 集群:%s 原因分析:%s",
						event.InvolvedObject.Kind, failedType, event.Name, event.Namespace, reason)
				}
				//调度失败
			}
		}
	}
}
