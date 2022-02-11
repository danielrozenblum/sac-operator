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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"

	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"

	"bitbucket.org/accezz-io/sac-operator/service"

	logger "sigs.k8s.io/controller-runtime/pkg/log"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SiteReconcile reconciles a Site object
type SiteReconcile struct {
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
func (r *SiteReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx).WithValues("sac-site spec", req.NamespacedName)

	siteCRD := &accessv1.Site{}
	connectorList := &corev1.PodList{}

	if err := r.Get(ctx, req.NamespacedName, siteCRD); err != nil {
		log.Error(err, "unable to fetch site")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.List(ctx, connectorList, client.InNamespace(siteCRD.Spec.ConnectorsNamespace), client.MatchingFields{podOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list pods Jobs")
		return ctrl.Result{}, err
	}

	log.WithValues("siteCRD", siteCRD, "connectorList", connectorList).Info("got triggered")

	switch {
	case siteCRD.Status.ID == nil:
		return r.CreateSite(ctx, siteCRD)
	default:
		return r.ReconcileConnectors(ctx, siteCRD, connectorList)
	}

}

func (r *SiteReconcile) CreateSite(ctx context.Context, siteCRD *accessv1.Site) (ctrl.Result, error) {
	log := logger.FromContext(ctx).WithValues("creating site", siteCRD.Name)

	siteModel, err := r.SiteConverter.ConvertToServiceModel(siteCRD)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.SiteService.Create(ctx, siteModel)
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

func (r *SiteReconcile) ReconcileConnectors(ctx context.Context, siteCRD *accessv1.Site, connectorList *corev1.PodList) (ctrl.Result, error) {
	_ = logger.FromContext(ctx).WithValues("ReconcileConnectors for site", siteCRD.Name)

	return ctrl.Result{}, nil
}

var (
	podOwnerKey = ".metadata.controller"
	apiGVStr    = accessv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *SiteReconcile) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Pod{}, podOwnerKey, func(rawObj client.Object) []string {
		// grab the pod object, extract the owner...
		pod := rawObj.(*corev1.Pod)
		owner := metav1.GetControllerOf(pod)
		if owner == nil {
			return nil
		}
		// ...make sure it's a Site...
		if owner.APIVersion != apiGVStr || owner.Kind != "Site" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&accessv1.Site{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
