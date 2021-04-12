package reconciler

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/6RiverSystems/operator-toolkit/apis"
)

// Reconciler incapsulates a bunch of helpful methods usually required by custom resource reconciler
type Reconciler struct {
	client    client.Client
	scheme    *runtime.Scheme
	recorder  record.EventRecorder
	finalizer string
	log       logr.Logger
}

// GetClient returns k8s API client
func (r *Reconciler) GetClient() client.Client { return r.client }

// GetScheme returns runtime scheme
func (r *Reconciler) GetScheme() *runtime.Scheme { return r.scheme }

// GetRecorder returns event recorder
func (r *Reconciler) GetRecorder() record.EventRecorder { return r.recorder }

// GetLogger returns reconciler logger
func (r *Reconciler) GetLogger() logr.Logger { return r.log }

func (r *Reconciler) loggerFor(obj Resource) logr.Logger {
	log := r.log.WithValues("Object.NamespacedName", path.Join(obj.GetNamespace(), obj.GetName()))
	if r.log.V(2).Enabled() {
		gvks, isUnversioned, err := r.scheme.ObjectKinds(obj)
		if err != nil {
			log.Error(err, "could not resolve object kind")
		} else if isUnversioned {
			log.Info("Unversioned resource")
		} else {
			log = log.WithValues("gvks", gvks)
		}
	}
	if r.log.V(3).Enabled() {
		if _, ok := obj.(*corev1.Secret); !ok {
			log = log.WithValues("object", obj)
		}
	}
	return log
}

// CreateResourceIfNotExists create a resource if it doesn't already exists. If the resource exists it is left untouched and the functin does not fails
// if owner is not nil, the owner field os set
// if obj namespace is "", the namespace field of the owner is assigned
func (r *Reconciler) CreateResourceIfNotExists(ctx context.Context, owner, obj Resource) error {
	if owner != nil {
		_ = controllerutil.SetControllerReference(owner, obj, r.GetScheme())
		if obj.GetNamespace() == "" {
			obj.SetNamespace(owner.GetNamespace())
		}
	}

	log := r.loggerFor(obj)
	log.V(2).Info("Creating resource if does not exist")

	err := r.GetClient().Create(ctx, obj)
	if err != nil && apierrors.IsAlreadyExists(err) {
		log.V(2).Info("Resource already exists")
	} else if err != nil {
		log.Error(err, "unable to create object")
		return err
	}
	return nil
}

// CreateOrUpdateResource creates a resource if it doesn't exist, and updates (overwrites it), if it exist
// if owner is not nil, the owner field os set
// if obj namespace is "", the namespace field of the owner is assigned
func (r *Reconciler) CreateOrUpdateResource(ctx context.Context, owner, obj Resource) error {
	if owner != nil {
		_ = controllerutil.SetControllerReference(owner, obj, r.GetScheme())
		if obj.GetNamespace() == "" {
			obj.SetNamespace(owner.GetNamespace())
		}
	}

	log := r.loggerFor(obj)

	obj2 := obj.DeepCopyObject()

	err := r.GetClient().Get(ctx, types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}, obj2)

	if apierrors.IsNotFound(err) {
		log.V(2).Info("Creating resource")
		err = r.GetClient().Create(ctx, obj)
		if err != nil {
			log.Error(err, "unable to create object")
			return err
		}
		return nil
	}
	anno := obj.GetAnnotations()
	if val := anno["control-tower/do-not-reconcile"]; val == "true" {
		return nil
	}
	if err == nil {
		obj3, ok := obj2.(metav1.Object)
		if !ok {
			return errors.New("unable to convert to metav1.Object")
		}
		obj.SetResourceVersion(obj3.GetResourceVersion())
		log.V(2).Info("Updating resource")
		err = r.GetClient().Update(ctx, obj)
		if err != nil {
			log.Error(err, "unable to update object")
			return err
		}
		return nil
	}
	r.log.Error(err, "unable to lookup object", "object", obj)
	return err
}

// CreateOrPatchResource creates a resource if it doesn't exist, and patches, if it exist
// if owner is not nil, the owner field os se
// if obj namespace is "", the namespace field of the owner is assigned
func (r *Reconciler) CreateOrPatchResource(ctx context.Context, owner, obj Resource, fpatch func(objFound, objNew runtime.Object) error) error {
	if owner != nil {
		_ = controllerutil.SetControllerReference(owner, obj, r.GetScheme())
		if obj.GetNamespace() == "" {
			obj.SetNamespace(owner.GetNamespace())
		}
	}

	found := obj.DeepCopyObject().(Resource)

	err := r.GetClient().Get(ctx, types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}, found)

	if apierrors.IsNotFound(err) {
		log := r.loggerFor(obj)
		log.V(2).Info("Creating resource")
		err = r.GetClient().Create(ctx, obj)
		if err != nil {
			log.Error(err, "unable to create object")
			return err
		}
		return nil
	}

	if err == nil {
		log := r.loggerFor(found)

		anno := obj.GetAnnotations()
		if val := anno["control-tower/do-not-reconcile"]; val == "true" {
			return nil
		}

		patch := client.MergeFrom(found.DeepCopyObject())
		if err := fpatch(found, obj); err != nil {
			log.Error(err, "failed to modify object")
			return err
		}

		log.V(2).Info("Patching resource")
		err = r.GetClient().Patch(ctx, found, patch)
		if err != nil {
			log.Error(err, "unable to patch object")
			return err
		}
		return nil
	}

	r.log.Error(err, "unable to lookup object", "object", obj)
	return err
}

// DeleteResourceIfExists deletes an existing resource. It doesn't fail if the resource does not exist
func (r *Reconciler) DeleteResourceIfExists(obj Resource) error {
	log := r.loggerFor(obj)
	log.V(2).Info("Removing resource", "object", obj)
	err := r.GetClient().Delete(context.TODO(), obj)
	if err != nil && !apierrors.IsNotFound(err) {
		log.Error(err, "unable to delete object ", "object", obj)
		return err
	}
	return nil
}

// DeleteResourcesIfExist operates like DeleteResources, but on an arrays of resources
func (r *Reconciler) DeleteResourcesIfExist(objs []Resource) error {
	for _, obj := range objs {
		err := r.DeleteResourceIfExists(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateResourcesIfNotExist operates as CreateResourceIfNotExists, but on an array of resources
func (r *Reconciler) CreateResourcesIfNotExist(ctx context.Context, owner Resource, objs []Resource) error {
	for _, obj := range objs {
		err := r.CreateResourceIfNotExists(ctx, owner, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdateResources operates as CreateOrUpdate, but on an array of resources
func (r *Reconciler) CreateOrUpdateResources(ctx context.Context, owner Resource, objs []Resource) error {
	for _, obj := range objs {
		err := r.CreateOrUpdateResource(ctx, owner, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

// ManageError will take care of the following:
// 1. generate a warning event attched to the passed object
// 2. set the status to failed if the object implements the readyStatusAware interface
// 3. return empty reconcile result with the passed error
func (r *Reconciler) ManageError(ctx context.Context, obj Resource, err error) (reconcile.Result, error) {
	r.recorder.Event(obj, "Warning", "ProcessingError", err.Error())
	log := r.loggerFor(obj)

	statusChanged := false

	// set condition if the error is of ConditionError type
	if conditionsStatusAware, ok := obj.(conditionsStatusAware); ok {
		var conditionErr *ConditionError
		if errors.As(err, &conditionErr) {
			c := conditionErr.Condition()
			log.V(1).Info("Setting status condition", "Condition", c)
			statusChanged = conditionsStatusAware.SetCondition(c)
			// unwrap error
			err = conditionErr.Err
		}
	}
	// Set readiness state if supported
	if readyStatusAware, ok := obj.(readyStatusAware); ok {
		log.V(2).Info("Setting readiness status to failed")
		readyStatusAware.SetReadyStatus(apis.FailedReadyStatus(err))
		statusChanged = true
	}

	// Update status if changed
	if statusChanged {
		if err := r.client.Status().Update(context.TODO(), obj); err != nil {
			log.Error(err, "Unable to update status")
			return reconcile.Result{RequeueAfter: time.Second}, nil
		}
	}

	// For not retriable errors set RequeueAfter to 6 hours (avg time of Event lifetime)
	var notRetriableErr *notRetriableError
	if errors.As(err, &notRetriableErr) {
		log.Error(notRetriableErr.issue, "Not retriable error")
		return reconcile.Result{RequeueAfter: 6 * time.Hour}, nil
	}

	// For retriable error, we have a recommended time interval to retry reconcile later.
	var retriableErr *retriableError
	if errors.As(err, &retriableErr) {
		log.Info("Retriable error. Mute an issue and reque", "Error", retriableErr.Error(), "RequeAfter", retriableErr.retryAfter)
		return reconcile.Result{RequeueAfter: retriableErr.retryAfter}, nil
	}

	return reconcile.Result{Requeue: true}, err
}

// ManageSuccess will update the status of the CR and return a successful reconcile result
func (r *Reconciler) ManageSuccess(ctx context.Context, obj Resource) (reconcile.Result, error) {
	if readyStatusAware, ok := obj.(readyStatusAware); ok {
		readyStatusAware.SetReadyStatus(apis.ReadyStatusOK())
	}
	if err := r.client.Status().Update(context.TODO(), obj); err != nil {
		r.loggerFor(obj).Error(err, "Unable to update status")
		return reconcile.Result{RequeueAfter: time.Second}, nil
	}
	return reconcile.Result{}, nil
}

// IsFinalized setup finalizer and updates the resource
func (r *Reconciler) IsFinalized(ctx context.Context, obj Resource, clean func() error) (bool, error) {
	log := r.loggerFor(obj).WithValues("finalizer", r.finalizer)

	if IsBeingDeleted(obj) {
		if !HasFinalizer(obj, r.finalizer) {
			return false, nil
		}

		if err := clean(); err != nil {
			log.Error(err, "Unable to cleanup finalizer")
			return false, err
		}

		log.V(2).Info("Removing finalizer")
		RemoveFinalizer(obj, r.finalizer)
		err := r.GetClient().Update(ctx, obj)
		if err != nil {
			log.Error(err, "Unable to remove finalizer")
		}
		return false, err
	}

	if !HasFinalizer(obj, r.finalizer) {
		r.log.V(2).Info("Adding finalizer")
		AddFinalizer(obj, r.finalizer)
		err := r.GetClient().Update(ctx, obj)
		if err != nil {
			log.Error(err, "Unable to add finalizer")
		}
		return false, err
	}

	return true, nil
}

// NewWithManager allocates new base reconciler using manager
func NewWithManager(mgr manager.Manager, controllerName string) *Reconciler {
	return New(
		mgr.GetClient(),
		mgr.GetScheme(),
		mgr.GetEventRecorderFor(controllerName),
		controllerName,
	)
}

// New allocates new base reconciler
func New(client client.Client, scheme *runtime.Scheme, recorder record.EventRecorder, controllerName string) *Reconciler {
	return &Reconciler{
		client:    client,
		scheme:    scheme,
		recorder:  recorder,
		finalizer: controllerName,
		log:       logf.Log.WithName(controllerName),
	}
}
