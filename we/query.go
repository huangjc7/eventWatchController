package we

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func (c *weController) queryPodRestart(pod []*v1.Pod, podName string) {
	for _, v := range pod {
		if v.Name == podName {
			fmt.Println("vvvvv"v.Spec.Containers[0].Resources.Limits)
		}
	}
}
