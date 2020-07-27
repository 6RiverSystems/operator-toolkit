package predicates

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/event"
	// logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ResourceGenerationOrFinalizerChangedPredicate passes event if either
// generation or finalizer is changed.
type ResourceGenerationOrFinalizerChangedPredicate struct {
	predicate.Funcs
}

// Update implements default UpdateEvent filter for validating resource version change
func (ResourceGenerationOrFinalizerChangedPredicate) Update(e event.UpdateEvent) bool {
	if e.MetaNew.GetGeneration() == e.MetaOld.GetGeneration() && reflect.DeepEqual(e.MetaNew.GetFinalizers(), e.MetaOld.GetFinalizers()) {
		return false
	}

	// log := logf.Log.WithName("ResourceGenerationOrFinalizerChangedPredicate")
	// log.V(4).Info("Triggering event",
	// 	"MetaNew.Generation", e.MetaNew.GetGeneration(),
	// 	"MetaOld.Generation", e.MetaOld.GetGeneration(),
	// 	"MetaNew.Finalizers", strings.Join(e.MetaNew.GetFinalizers(), ","),
	// 	"MetaOld.Finalizers", strings.Join(e.MetaOld.GetFinalizers(), ","),
	// )

	return true
}
