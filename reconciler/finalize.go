package reconciler

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsBeingDeleted returns whether this object has been requested to be deleted
func IsBeingDeleted(obj metav1.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// HasFinalizer returns whether this object has the passed finalizer
func HasFinalizer(obj metav1.Object, finalizer string) bool {
	for _, fin := range obj.GetFinalizers() {
		if fin == finalizer {
			return true
		}
	}
	return false
}

// AddFinalizer adds the passed finalizer this object
func AddFinalizer(obj metav1.Object, finalizer string) {
	controllerutil.AddFinalizer(obj, finalizer)
}

// RemoveFinalizer removes the passed finalizer from object
func RemoveFinalizer(obj metav1.Object, finalizer string) {
	controllerutil.RemoveFinalizer(obj, finalizer)
}
