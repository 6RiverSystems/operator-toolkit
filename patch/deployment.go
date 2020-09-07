package patch

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ReplaceDeploymentSpec replaces only Spec field of the deployment
func ReplaceDeploymentSpec(objFound, objNew runtime.Object) error {
	var (
		ok       bool
		dst, src *appsv1.Deployment
	)
	if dst, ok = objFound.(*appsv1.Deployment); !ok {
		return fmt.Errorf("found object is not of deployment type: %+v", objFound)
	}
	if src, ok = objNew.(*appsv1.Deployment); !ok {
		return fmt.Errorf("new object is not of deployment type: %+v", objNew)
	}
	dst.Spec = src.Spec
	return nil
}
