package connector_deployer

import (
	"context"
	"fmt"
	"time"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/runtime"

	"bitbucket.org/accezz-io/sac-operator/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const annotationPrefix = "access.secure-access-cloud.symantec.com"

type KubernetesImpl struct {
	client.Client
	Scheme              *runtime.Scheme
	SiteNamespace       string
	ConnectorsNamespace string
}

func NewKubernetesImpl(client client.Client, scheme *runtime.Scheme) *KubernetesImpl {
	return &KubernetesImpl{Client: client, Scheme: scheme}
}

func (k *KubernetesImpl) CreateConnector(ctx context.Context, inputs *CreateConnectorInput) (*CreateConnctorOutput, error) {

	site, err := k.getSite(ctx, inputs.SiteName)
	if err != nil {
		return nil, err
	}

	pod, err := k.getConnectorPodForSite(inputs, site)
	if err != nil {
		return nil, err
	}
	err = k.Create(ctx, pod)
	if err != nil {
		return nil, err
	}

	uid, err := utils.FromUIDType(&pod.UID)
	if err != nil {
		return nil, err
	}

	return &CreateConnctorOutput{DeploymentID: uid}, nil
}

var (
	podOwnerKey = ".metadata.controller"
)

func (k *KubernetesImpl) GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error) {

	site, err := k.getSite(ctx, siteName)
	if err != nil {
		return []Connector{}, err
	}

	connectorList := &corev1.PodList{}
	if err := k.List(ctx, connectorList, client.InNamespace(site.Spec.ConnectorsNamespace), client.MatchingFields{podOwnerKey: site.Name}); err != nil {
		return []Connector{}, err
	}

	connectors := []Connector{}
	for i := range connectorList.Items {
		uid, err := utils.FromUIDType(&connectorList.Items[i].UID)
		if err != nil {
			return []Connector{}, err
		}
		switch connectorList.Items[i].Status.Phase {
		case corev1.PodRunning:
			connectors = append(connectors, Connector{
				ID:     uid,
				Status: OKConnectorStatus,
			})
		case corev1.PodFailed:
			connectors = append(connectors, Connector{
				ID:     uid,
				Status: RecreateConnectorStatus,
			})
		case corev1.PodPending:
			if time.Since(connectorList.Items[i].GetCreationTimestamp().Time) > 2*time.Minute {
				connectors = append(connectors, Connector{
					ID:     uid,
					Status: RecreateConnectorStatus,
				})
				continue
			}
			connectors = append(connectors, Connector{
				ID:     uid,
				Status: PendingConnectorStatus,
			})
		}
	}

	return connectors, nil
}

func (k *KubernetesImpl) getConnectorPodForSite(inputs *CreateConnectorInput, site *accessv1.Site) (*corev1.Pod, error) {

	podEnvVar := []corev1.EnvVar{}

	for key, val := range inputs.EnvironmentVars {
		podEnvVar = append(podEnvVar, corev1.EnvVar{
			Name:  key,
			Value: val,
		})
	}

	podName := inputs.Name
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    make(map[string]string),
			Namespace: k.ConnectorsNamespace,
			Name:      podName,
			Annotations: map[string]string{
				fmt.Sprintf("%s/%s", annotationPrefix, "connector"): inputs.ConnectorID.String(),
				fmt.Sprintf("%s/%s", annotationPrefix, "site"):      inputs.SiteName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "connector",
				Image: inputs.Image,
				Env:   podEnvVar,
			}},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:  utils.FromInt64(1000),
				RunAsGroup: utils.FromInt64(1000),
			},
		},
	}

	err := ctrl.SetControllerReference(site, pod, k.Scheme)
	if err != nil {
		return nil, err
	}

	return pod, nil

}

func (k *KubernetesImpl) getSite(ctx context.Context, siteName string) (*accessv1.Site, error) {
	site := &accessv1.Site{}
	if err := k.Get(ctx, client.ObjectKey{
		Namespace: k.SiteNamespace,
		Name:      siteName,
	}, site); err != nil {
		return nil, err
	}
	return site, nil
}
