package connector_deployer

import (
	"context"
	"fmt"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/runtime"

	"bitbucket.org/accezz-io/sac-operator/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

const annotationPrefix = "access.secure-access-cloud.symantec.com"

type KubernetesImpl struct {
	client.Client
	Scheme *runtime.Scheme
}

func NewKubernetesImpl(client client.Client, scheme *runtime.Scheme) *KubernetesImpl {
	return &KubernetesImpl{Client: client, Scheme: scheme}
}

func (k *KubernetesImpl) Deploy(ctx context.Context, inputs *DeployConnectorInput) (*DeployConnectorOutput, error) {

	site := &accessv1.Site{}
	if err := k.Get(context.Background(), client.ObjectKey{
		Namespace: inputs.SiteNamespace,
		Name:      inputs.SiteName,
	}, site); err != nil {
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

	return &DeployConnectorOutput{DeploymentID: uid}, nil
}

func (k *KubernetesImpl) getConnectorPodForSite(inputs *DeployConnectorInput, site *accessv1.Site) (*corev1.Pod, error) {

	podEnvVar := []corev1.EnvVar{}

	for key, val := range inputs.EnvironmentVars {
		podEnvVar = append(podEnvVar, corev1.EnvVar{
			Name:  key,
			Value: val,
		})
	}

	podName := fmt.Sprintf("%s-%s-%s", "connector", inputs.SiteName, rand.String(4))
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    make(map[string]string),
			Namespace: inputs.Namespace,
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
