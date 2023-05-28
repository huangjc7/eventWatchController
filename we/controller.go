package we

import v1 "k8s.io/api/core/v1"

type weController struct {
	Kubeconfig string
	Events     []*v1.Event
	StopCh     chan struct{}
	Pod        []*v1.Pod
}

var We weController
