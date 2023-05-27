package we

import v1 "k8s.io/api/core/v1"

type weController struct {
	Kubeconfig string
	events     []*v1.Event
	StopCh     chan struct{}
}

var We weController
