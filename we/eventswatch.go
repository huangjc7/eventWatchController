package we

import (
	"log"
	"regexp"
)

func (c *weController) Run(regexFailed, failedType, reason string) {
	for _, event := range c.Events {
		if event.Count >= 3 {
			// 启动失败
			if event.Reason == "BackOff" || event.Reason == "Failed" {
				if matched, _ := regexp.Match(regexFailed, []byte(event.Message)); matched {
					log.Printf("发现异常! 资源类型:%s 异常种类:%s 应用名字:%s 所在租户:%s 原因分析:%s",
						event.InvolvedObject.Kind, failedType, event.InvolvedObject.Name, event.Namespace, reason)
					if failedType == "重启失败" {
						c.queryPodRestart(c.Pod, event.InvolvedObject.Name)
					}
					//} else if failedType == "资源不足" {
					//	c.queryResource
					//} else if failedType == "存储卷" {
					//	c.queryVolume
					//} else if failedType == "创建失败" {
					//	c.queryPodCreateFailed
					//}
				}
				//if regexFailed
			}
		}
	}
}
