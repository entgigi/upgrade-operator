package common

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BaseK8sStructure struct {
	client.Client
	Log logr.Logger
}
