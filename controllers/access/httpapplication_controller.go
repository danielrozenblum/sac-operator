/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package access

import (
	"context"
	"errors"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"bitbucket.org/accezz-io/sac-operator/utils/typederror"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"bitbucket.org/accezz-io/sac-operator/service"

	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"

	"github.com/go-logr/logr"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HttpApplicationReconciler reconciles a HttpApplication object
type HttpApplicationReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	ApplicationService service.ApplicationService
	ConverterToModel   *converter.HttpApplicationTypeConverter
	Log                logr.Logger
}

//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=httpapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=httpapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=httpapplications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HttpApplication object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *HttpApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	application := &accessv1.HttpApplication{}

	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		r.Log.Error(err, "unable to fetch application")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	model, err := r.ConverterToModel.ConvertToModel(application)
	if err != nil {
		r.Log.Error(err, "convert to service model")
		return ctrl.Result{}, nil
	}
	output, err := r.ApplicationService.Reconcile(ctx, model)
	if !controllerutil.ContainsFinalizer(application, applicationFinalizerName) && output.SACApplicationID != "" {
		controllerutil.AddFinalizer(application, applicationFinalizerName)
		if err := r.Update(ctx, application); err != nil {
			r.Log.WithValues("application", application.Name).Info("failed to add finalizer")
			return ctrl.Result{}, err
		}
	}
	return r.handleReconcilerReturn(ctx, application, output, err)

}

// SetupWithManager sets up the controller with the Manager.
func (r *HttpApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1.HttpApplication{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *HttpApplicationReconciler) handleReconcilerReturn(ctx context.Context, application *accessv1.HttpApplication, output *service.ApplicationReconcileOutput, reconcileError error) (ctrl.Result, error) {

	if errors.Is(reconcileError, typederror.UnrecoverableError) {
		r.Log.WithValues("application", application.Name).Error(reconcileError, "got unrecoverable error, giving up...")
		return ctrl.Result{Requeue: false}, nil
	}

	if output.Deleted {
		controllerutil.RemoveFinalizer(application, applicationFinalizerName)
		if err := r.Update(ctx, application); err != nil {
			r.Log.WithValues("application", application.Name).Error(err, "failed to remove Finalizer from application")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	application.Status = r.ConverterToModel.ConvertFromServiceOutput(output)

	if reconcileError != nil {
		r.Log.WithValues("application", application.Name).Error(reconcileError, "failed to reconcile, trying to update last known status")
	}

	err := r.Status().Update(ctx, application)
	if err != nil {
		r.Log.WithValues("application", application.Name).Error(reconcileError, "failed to update application status, retrying in 5 seconds")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	if reconcileError != nil {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, reconcileError
	}

	return ctrl.Result{Requeue: false}, nil
}
