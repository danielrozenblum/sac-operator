package repository

import (
	"context"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type K8sImpl struct {
	client.Client
	Namespace string
}

func NewK8sImpl(client client.Client, namespace string) *K8sImpl {
	return &K8sImpl{Client: client, Namespace: namespace}
}

func (k *K8sImpl) UpdateNewSite(ctx context.Context, siteName, id string) error {
	site := &accessv1.Site{}

	if err := k.Get(ctx, client.ObjectKey{Namespace: k.Namespace, Name: siteName}, site); err != nil {
		return err
	}

	site.Status.ID = id

	if err := k.Status().Update(ctx, site); err != nil {
		return err
	}

	if !controllerutil.ContainsFinalizer(site, siteFinalizerName) {
		controllerutil.AddFinalizer(site, siteFinalizerName)
		if err := k.Update(ctx, site); err != nil {
			return err
		}
	}
	return nil

}

func (k *K8sImpl) UpdateDeleteSite(ctx context.Context, siteName string) error {
	site := &accessv1.Site{}

	if err := k.Get(ctx, client.ObjectKey{Namespace: k.Namespace, Name: siteName}, site); err != nil {
		return err
	}

	controllerutil.RemoveFinalizer(site, siteFinalizerName)
	if err := k.Update(ctx, site); err != nil {
		return err
	}
	return nil

}
