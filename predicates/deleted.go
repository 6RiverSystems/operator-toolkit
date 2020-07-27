package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ResourceDeletedPredicate allows an event to pass only when resource is deleted
type ResourceDeletedPredicate struct {
	predicate.Funcs
}

// Create events are ignored
func (ResourceDeletedPredicate) Create(event.CreateEvent) bool { return false }

// Update events are ignored
func (ResourceDeletedPredicate) Update(event.UpdateEvent) bool { return false }

// Generic events are ignored
func (ResourceDeletedPredicate) Generic(event.GenericEvent) bool { return false }
