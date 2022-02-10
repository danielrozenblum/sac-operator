package model

import (
	"fmt"

	"bitbucket.org/accezz-io/sac-operator/utils"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type Connector struct {
	ConnectorID           *uuid.UUID
	ConnectorDeploymentID *uuid.UUID
	Version               string
	Name                  string
	ConnectorStatus       string
}

type Site struct {
	Name                string
	ID                  *uuid.UUID
	TenantIdentifier    string
	NumberOfConnectors  int
	ConnectorsNamespace string
	EndpointURL         string
	Connectors          []Connector
}

const annotationPrefix = "access.secure-access-cloud.symantec.com"
const connectorImage = "luminate/connector"

func (s *Site) getConnectorPodForSite(version, connectorID, otp string) *corev1.Pod {

	podName := fmt.Sprintf("%s-%s-%s", "connector", s.Name, rand.String(4))
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    make(map[string]string),
			Namespace: s.ConnectorsNamespace,
			Name:      podName,
			Annotations: map[string]string{
				fmt.Sprintf("%s/%s", annotationPrefix, "connector"): connectorID,
				fmt.Sprintf("%s/%s", annotationPrefix, "site"):      s.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "connector",
				Image: fmt.Sprintf("%s/%s", connectorImage, version),
				Env: []corev1.EnvVar{
					{
						Name:  "HTTPS_SKIP_CERT_VERIFY",
						Value: "true",
					},
					{
						Name:  "OTP",
						Value: otp,
					},
					{
						Name:  "TENANT_IDENTIFIER",
						Value: s.TenantIdentifier,
					},
					{
						Name:  "ENDPOINT_URL",
						Value: s.EndpointURL,
					},
				},
			}},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:  utils.FromInt64(1000),
				RunAsGroup: utils.FromInt64(1000),
			},
		},
	}

	return pod

}
