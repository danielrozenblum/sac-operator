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
	"time"

	corev1 "k8s.io/api/core/v1"

	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"

	"bitbucket.org/accezz-io/sac-operator/service"

	logger "sigs.k8s.io/controller-runtime/pkg/log"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SiteReconciler reconciles a Site object
type SiteReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	SiteService   service.SiteService
	SiteConverter *converter.SiteConverter
}

//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=sites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=sites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=access.secure-access-cloud.symantec.com,resources=sites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Site object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx).WithValues("sac-site spec", req.NamespacedName)

	siteCRD := &accessv1.Site{}

	if err := r.Get(ctx, req.NamespacedName, siteCRD); err != nil {
		log.Error(err, "unable to fetch site")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	switch {
	case siteCRD.Status.ID == nil:
		siteModel, err := r.SiteConverter.ConvertToServiceModel(siteCRD)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.SiteService.Create(ctx, siteModel, siteCRD)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.SiteConverter.UpdateStatus(siteModel, &siteCRD.Status)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.Status().Update(ctx, siteCRD); err != nil {
			log.Error(err, "unable to update siteCRD status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: 30 * time.Second,
		}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1.Site{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
