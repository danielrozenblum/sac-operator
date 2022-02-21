package connector_deployer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	ctrl "sigs.k8s.io/controller-runtime"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/runtime"

	"bitbucket.org/accezz-io/sac-operator/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AnnotationPrefix = "access.secure-access-cloud.symantec.com"

type KubernetesImpl struct {
	client.Client
	Scheme              *runtime.Scheme
	SiteNamespace       string
	ConnectorsNamespace string
	podOwnerKey         string
	log                 logr.Logger
}

func NewKubernetesImpl(client client.Client, scheme *runtime.Scheme, podOwnerKey string, log logr.Logger) *KubernetesImpl {
	return &KubernetesImpl{Client: client, Scheme: scheme, podOwnerKey: podOwnerKey, log: log}
}

func (k *KubernetesImpl) CreateConnector(ctx context.Context, inputs *CreateConnectorInput) (string, error) {

	site, err := k.getSite(ctx, inputs.SiteName)
	if err != nil {
		return "", err
	}

	pod := k.getConnectorPodForSite(inputs, site)
	k.log.WithValues("pod", pod.Name).Info("creating connector in k8s")
	err = k.Create(ctx, pod)
	if err != nil {
		return "", err
	}

	// TODO implement using watch - https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/source/example_test.go
	count := 5
	for i := 0; i < count; i-- {
		if count == 0 {
			return "", fmt.Errorf("onnector failed to reach stable state")
		}
		time.Sleep(time.Second * 10)
		err = k.Get(ctx, client.ObjectKey{Namespace: pod.Namespace, Name: pod.Name}, pod)
		if err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return "", err
		}
		if pod.Status.Phase == corev1.PodRunning && pod.Status.ContainerStatuses[0].Ready {
			return pod.ObjectMeta.Name, err
		}
		k.log.WithValues("pod", pod.Name).Info("waiting for pod and container to be in running state")
		count--
	}

	return pod.Name, nil
}

func (k *KubernetesImpl) DeleteConnector(ctx context.Context, name string) error {

	podToDelete := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: k.ConnectorsNamespace,
			Name:      name,
		},
		Spec:   corev1.PodSpec{},
		Status: corev1.PodStatus{},
	}

	err := k.Delete(ctx, podToDelete)
	count := 5
	for i := 0; i < count; i-- {
		if count == 0 {
			return fmt.Errorf("failed to valied connector deletion from cluster")
		}
		time.Sleep(time.Second * 10)
		err = k.Get(ctx, client.ObjectKey{Namespace: k.ConnectorsNamespace, Name: name}, podToDelete)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to valied connector deletion from cluster")
		}
		k.log.WithValues("pod", name).Info("waiting for pod to be deleted")
		count--
	}

	return client.IgnoreNotFound(err)

}

func (k *KubernetesImpl) GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error) {

	site, err := k.getSite(ctx, siteName)
	if err != nil {
		return []Connector{}, err
	}

	connectorList := &corev1.PodList{}
	if err := k.List(ctx, connectorList, client.InNamespace(site.Spec.ConnectorsNamespace),
		client.MatchingFields{k.podOwnerKey: site.Name}); err != nil {
		return []Connector{}, err
	}

	connectors := []Connector{}
	for i := range connectorList.Items {
		sacID := connectorList.Items[i].GetAnnotations()[k.connectorAnnotationKey()]
		connector := Connector{}
		connector.DeploymentName = connectorList.Items[i].GetName()
		connector.SACID = sacID
		connector.CreatedTimeStamp = connectorList.Items[i].GetCreationTimestamp().Time
		switch connectorList.Items[i].Status.Phase {
		case corev1.PodRunning:
			if connectorList.Items[i].Status.ContainerStatuses[0].Ready {
				connector.Status = OKConnectorStatus
				break
			}
			connector.Status = ToDeleteConnectorStatus
		case corev1.PodFailed:
			connector.Status = ToDeleteConnectorStatus
		case corev1.PodPending:
			if time.Since(connectorList.Items[i].GetCreationTimestamp().Time) > 2*time.Minute {
				connector.Status = ToDeleteConnectorStatus
				break
			}
			connector.Status = OKConnectorStatus
		}
		connectors = append(connectors, connector)
	}

	return connectors, nil
}

func (k *KubernetesImpl) getConnectorPodForSite(inputs *CreateConnectorInput, site *accessv1.Site) *corev1.Pod {

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
				fmt.Sprintf("%s/%s", AnnotationPrefix, "connector"): inputs.ConnectorID,
				fmt.Sprintf("%s/%s", AnnotationPrefix, "site"):      inputs.SiteName,
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

	_ = ctrl.SetControllerReference(site, pod, k.Scheme)

	return pod

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

func (k *KubernetesImpl) connectorAnnotationKey() string {
	return fmt.Sprintf("%s/%s", AnnotationPrefix, "connector")
}
