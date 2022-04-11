package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {

	// The service name
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// The service namespace (default is the application's namespace)
	// +optional
	Namespace string `json:"namespace"`

	// The port name that will be exposed by this application.
	// +kubebuilder:validation:Required
	Port string `json:"port"`

	// Protocol Schema (default is based on port and application type)
	// +optional
	Schema string `json:"schema,omitempty"`
}

type CommonApplicationParams struct {

	// The site to bind this application. The site should be an existing Site in your Secure Access Cloud tenant
	SiteName string `json:"site"`

	// A list of access-policies names to enforce on this application.
	// +optional
	AccessPoliciesNames []string `json:"access_policies,omitempty"`

	// A list of activity-policies names to enforce on this application.
	// +optional
	ActivityPoliciesNames []string `json:"activity_policies,omitempty"`

	// +kubebuilder:default=true
	IsVisible *bool `json:"is_visible,omitempty"`

	// +kubebuilder:default=false
	IsNotificationEnabled *bool `json:"is_notification_enabled,omitempty"`

	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`
}

type CommonApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The application-id in Secure-Access-Cloud
	// +optional
	Id string `json:"id,omitempty"`

	// Information when was the last time the application was successfully modified by the operator.
	// +optional
	ModifiedOn metav1.Time `json:"modifiedOn,omitempty"`
}
