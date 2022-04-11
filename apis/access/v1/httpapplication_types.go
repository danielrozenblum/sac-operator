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

package v1

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HttpApplicationSpec defines the desired state of HttpApplication
type HttpApplicationSpec struct {
	CommonApplicationParams `json:",inline"`
	// SubType of the application. Valid values are: HTTP_LUMINATE_DOMAIN, HTTP_CUSTOM_DOMAIN, HTTP_WILDCARD_DOMAIN
	// (default is HTTP_LUMINATE_DOMAIN)
	// +kubebuilder:validation:Enum=HTTP_LUMINATE_DOMAIN;HTTP_CUSTOM_DOMAIN;HTTP_WILDCARD_DOMAIN
	// +kubebuilder:default=HTTP_LUMINATE_DOMAIN
	SubType model.ApplicationSubType `json:"sub_type,omitempty"`

	Service `json:"service"`
	// +optional
	*HttpConnectionSettings `json:"connection_settings"`
	// +optional
	*HttpLinkTranslationSettings `json:"link_translation_settings"`
	// +optional
	*HttpRequestCustomizationSettings `json:"request_customization_settings"`
}

//// HttpApplicationStatus defines the observed state of HttpApplication
//type HttpApplicationStatus struct {
//	CommonApplicationStatus
//}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HttpApplication is the Schema for the httpapplications API
type HttpApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpApplicationSpec     `json:"spec,omitempty"`
	Status CommonApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HttpApplicationList contains a list of HttpApplication
type HttpApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpApplication{}, &HttpApplicationList{})
}

type HttpConnectionSettings struct {
	// +kubebuilder:validation:Optional
	SubDomain string `json:"subdomain,omitempty"`
	// +kubebuilder:validation:Optional
	CustomExternalAddress string `json:"custom_external_address,omitempty"`
	// +kubebuilder:validation:Optional
	CustomRootPath string `json:"custom_root_path,omitempty"`
	// +kubebuilder:validation:Optional
	HealthUrl string `json:"health_url,omitempty"`
	// +kubebuilder:validation:Optional
	HealthMethod string `json:"health_method,omitempty"`
	// +kubebuilder:validation:Optional
	CustomSSLCertificate string `json:"custom_ssl_certificate,omitempty"`
	// +kubebuilder:validation:Optional
	WildcardPrivateKey string `json:"wildcard_private_key,omitempty"`
}

type HttpLinkTranslationSettings struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	IsDefaultContentRewriteRulesEnabled *bool `json:"is_default_content_rewrite_rules_enabled"`
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	IsDefaultHeaderRewriteRulesEnabled *bool `json:"is_default_header_rewrite_rules_enabled"`
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	UseExternalAddressForHostAndSni *bool `json:"use_external_address_for_host_and_sni"`
	// +kubebuilder:validation:Optional
	LinkedApplications []string `json:"linked_applications"`
}

type HttpRequestCustomizationSettings struct {
	HeaderCustomization map[string]string `json:"header_customization"`
}
