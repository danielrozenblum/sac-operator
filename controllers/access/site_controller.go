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

	"bitbucket.org/accezz-io/sac-operator/service/repository"

	"github.com/go-logr/logr"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/service/sac"

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
	Scheme                    *runtime.Scheme
	SiteConverter             *converter.SiteConverter
	SecureAccessCloudSettings *sac.SecureAccessCloudSettings
	Log                       logr.Logger
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

	if err := r.Get(ctx, req.NamespacedName, siteCRD); err != nil {
		log.Error(err, "unable to fetch site")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	siteModel := r.SiteConverter.ConvertToServiceModel(siteCRD)

	serviceImpl := r.ServiceFacotry(ctx, siteCRD)

	output, err := serviceImpl.Reconcile(ctx, siteModel)
	siteCRD.Status = r.SiteConverter.ConvertFromServiceOutput(output)
	if err != nil {
		log.Error(err, "got error from reconcile, trying to update last know status")
		err = r.Status().Update(ctx, siteCRD)
		if err != nil {
			log.Error(err, "go error from reconcile, trying to update last know status")
		}
	}
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

func (r *SiteReconcile) ServiceFacotry(ctx context.Context, site *accessv1.Site) *service.SiteServiceImpl {

	sacClient := sac.NewSecureAccessCloudClientImpl(r.SecureAccessCloudSettings)
	k8sClients := connector_deployer.NewKubernetesImpl(r.Client, r.Scheme, podOwnerKey)
	k8sClients.ConnectorsNamespace = site.Spec.ConnectorsNamespace
	k8sClients.SiteNamespace = site.Namespace
	repo := repository.NewK8sImpl(r.Client, site.Namespace)
	return service.NewSiteServiceImpl(sacClient, k8sClients, r.Log, repo)

}
