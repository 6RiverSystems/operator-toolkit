package finalize

import (
	"context"

	"github.com/6RiverSystems/operator-toolkit/strslice"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Finalized interface is an interface object need to implement to be finalized
type Finalized interface {
	runtime.Object
	GetDeletionTimestamp() *metav1.Time
	GetFinalizers() []string
	SetFinalizers([]string)
}

// Result is finalization result
type Result struct {
	IsInitialized bool // is set to true if missing finalizer is added to resource
	IsFinalized   bool // is set to true if finalization completes

}

// Finalizer is a controller to handle CR finalization process
type Finalizer struct {
	name   string // finalizer name
	client client.Client
}

// NewFinalizer construct a new finalizer
func NewFinalizer(name string, client client.Client) *Finalizer {
	return &Finalizer{name: name, client: client}
}

// Finalize performs object finalization.
// Returns Result{IsFinalized: true} if object was finalized, otherwise Result{IsFinalized: false} - indicating that
// an object is not a subject for finalization action or finalization failed with error.
func (f *Finalizer) Finalize(reqLogger logr.Logger, m Finalized, clean func() error) (Result, error) {
	// Check if the instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isMarkedToBeDeleted := m.GetDeletionTimestamp() != nil
	// Check if resource contains finalizer
	containsFinalizer := strslice.Contains(m.GetFinalizers(), f.name)

	if isMarkedToBeDeleted {
		if containsFinalizer {
			// Run finalization logic.
			// If the finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := clean(); err != nil {
				return Result{}, err
			}

			// Remove finalizer. Once all finalizers have been removed, the object will be deleted.
			if err := f.removeFinalizer(reqLogger, m); err != nil {
				return Result{}, err
			}
		}
		return Result{IsFinalized: true}, nil
	}

	// add finalizer if it's missing
	if !containsFinalizer {
		if err := f.addFinalizer(reqLogger, m); err != nil {
			return Result{}, err
		}
		return Result{IsInitialized: true}, nil
	}
	return Result{}, nil
}

func (f *Finalizer) removeFinalizer(reqLogger logr.Logger, m Finalized) error {
	reqLogger.Info("Removing Finalizer")
	m.SetFinalizers(strslice.Remove(m.GetFinalizers(), f.name))
	err := f.client.Update(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}

func (f *Finalizer) addFinalizer(reqLogger logr.Logger, m Finalized) error {
	reqLogger.Info("Adding Finalizer")
	m.SetFinalizers(append(m.GetFinalizers(), f.name))
	err := f.client.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update resource object with finalizer")
		return err
	}
	return nil
}
