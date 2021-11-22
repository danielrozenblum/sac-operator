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
	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service"
	"bitbucket.org/accezz-io/sac-operator/utils"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logger "sigs.k8s.io/controller-runtime/pkg/log"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	ApplicationService   service.ApplicationService
	ApplicationConverter *converter.ApplicationTypeConverter
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
	log := logger.FromContext(ctx).WithValues("sac-application spec", req.NamespacedName)

	// 1. Load the application by name (namespaced name) from the Kubernetes Cluster (both spec & status).
	var application accessv1.Application
	var updatedApplicationOnSac *model.Application

	if err := r.Get(ctx, req.NamespacedName, &application); err != nil {
		log.Error(err, "unable to fetch application")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		// TODO should delete application here? if the application not found does it means it got deleted?
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Compare to the application configured on Secure-Access-Cloud:
	// 	2.1. In case not found, create the application (application-id not found on the status)
	if application.Status.Id == nil {
		var err error
		applicationModel := r.ApplicationConverter.ConvertToModel(application)
		updatedApplicationOnSac, err = r.ApplicationService.Create(ctx, applicationModel)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		//	2.2. In case found, compare the application configured on Secure-Access-Cloud to the desired state in the spec
		applicationId, err := utils.FromUIDType(application.Status.Id)
		if err != nil {
			return ctrl.Result{}, err
		}

		applicationModel := r.ApplicationConverter.ConvertToModel(application)
		updatedApplicationOnSac, err = r.ApplicationService.Update(ctx, applicationId, applicationModel)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// 3. Update the application status
	err := r.updateApplicationState(ctx, log, &application, updatedApplicationOnSac)
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: fetch the exposed service reference in order to updates on the service (port, service deleted)

	log.V(1).Info(
		"application reconciled successfully",
		"id", updatedApplicationOnSac.ID,
		"name", updatedApplicationOnSac.Name,
	)

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) updateApplicationState(
	ctx context.Context,
	log logr.Logger,
	application *accessv1.Application,
	updatedApplicationOnSac *model.Application,
) error {
	application.Status.Id = utils.FromUUID(updatedApplicationOnSac.ID)

	if err := r.Status().Update(ctx, application); err != nil {
		log.Error(err, "unable to update Application status")
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1.Application{}).
		Complete(r)
}
