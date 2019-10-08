package webhook

import (
	"github.com/jcrossley3/k-s-o-minikube/pkg/webhook/knativeserving"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, knativeserving.Add)
}
