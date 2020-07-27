package reconciler

import (
	"github.com/6RiverSystems/operator-toolkit/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource represents a kubernetes Resource
type Resource interface {
	metav1.Object
	runtime.Object
}

type readyStatusAware interface {
	GetReadyStatus() apis.ReadyStatus
	SetReadyStatus(status apis.ReadyStatus)
}

type conditionsStatusAware interface {
	SetCondition(apis.Condition) bool
}
