package patch

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ReplaceServicePorts replaces only Ports field of the service
func ReplaceServicePorts(objFound, objNew runtime.Object) error {
	var (
		ok       bool
		dst, src *corev1.Service
	)
	if dst, ok = objFound.(*corev1.Service); !ok {
		return fmt.Errorf("found object is not of service type: %+v", objFound)
	}
	if src, ok = objNew.(*corev1.Service); !ok {
		return fmt.Errorf("new object is not of service type: %+v", objNew)
	}
	dst.Spec.Ports = src.Spec.Ports
	return nil
}
