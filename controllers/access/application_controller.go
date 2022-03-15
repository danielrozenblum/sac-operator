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
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"bitbucket.org/accezz-io/sac-operator/utils/typederror"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"
	"bitbucket.org/accezz-io/sac-operator/service"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	ApplicationService   service.ApplicationService
	ApplicationConverter *converter.ApplicationTypeConverter
	Log                  logr.Logger
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1.Application{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to move the current state of the cluster closer to the desired state.
// This function should compare the state specified by the Application object against the actual cluster state, and then perform operations to
// make the cluster state reflect the state specified by the user.
//
// Reconcile implementations compare the state specified in an object by a user against the actual cluster state, and then perform operations
// to make the actual cluster state reflect the state specified by the user.
//
// The Controller will requeue the Request to be processed again if an error is non-nil or Result.Requeue is true,
// otherwise upon completion it will remove the work from the queue.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	application := &accessv1.Application{}
	//var updatedApplicationOnSac *model.Application

	r.Log.Info(fmt.Sprintf("got reconcile request for %+v", req.NamespacedName))

	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		r.Log.Error(err, "unable to fetch application")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.ApplicationConverter.Validate(application); err != nil {
		r.Log.Error(err, "validation error")
		return ctrl.Result{}, nil
	}

	model := r.ApplicationConverter.ConvertToModel(application)
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

func (r *ApplicationReconciler) handleReconcilerReturn(ctx context.Context, application *accessv1.Application, output *service.ApplicationReconcileOutput, reconcileError error) (ctrl.Result, error) {

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

	application.Status = r.ApplicationConverter.ConvertFromServiceOutput(output)

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
