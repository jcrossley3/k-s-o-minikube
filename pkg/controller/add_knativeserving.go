package controller

import (
	"github.com/jcrossley3/k-s-o-minikube/pkg/controller/knativeserving"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, knativeserving.Add)
}
