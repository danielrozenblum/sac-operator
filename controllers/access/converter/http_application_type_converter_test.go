package converter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/accezz-io/sac-operator/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/stretchr/testify/assert"
)

func TestHttpApplicationTypeConverter_ConvertToModel(t *testing.T) {
	type args struct {
		application *accessv1.HttpApplication
	}
	tests := []struct {
		name        string
		args        args
		expected    *model.Application
		errorOutput error
	}{
		{
			name: "explicit flow",
			args: args{
				application: &accessv1.HttpApplication{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-application",
					},
					Spec: accessv1.HttpApplicationSpec{
						SubType: "HTTP_CUSTOM_DOMAIN",
						CommonApplicationParams: accessv1.CommonApplicationParams{
							SiteName:              "my-site",
							AccessPoliciesNames:   []string{"access-policy"},
							ActivityPoliciesNames: []string{"activity-policy"},
							IsVisible:             utils.Convert_bool_To_Pointer_bool(true),
							IsNotificationEnabled: utils.Convert_bool_To_Pointer_bool(true),
							Enabled:               utils.Convert_bool_To_Pointer_bool(true),
						},
						Service: accessv1.Service{
							Name:      "my-service-name",
							Namespace: "service-namespace",
							Port:      "80",
							Schema:    "http",
						},
						HttpConnectionSettings: &accessv1.HttpConnectionSettings{
							SubDomain:             "subDomain",
							CustomExternalAddress: "customExternalAddress",
							CustomRootPath:        "/root/path",
							HealthUrl:             "/health",
							HealthMethod:          "GET",
							CustomSSLCertificate:  "",
							WildcardPrivateKey:    "",
						},
						HttpLinkTranslationSettings: &accessv1.HttpLinkTranslationSettings{
							IsDefaultContentRewriteRulesEnabled: utils.Convert_bool_To_Pointer_bool(false),
							IsDefaultHeaderRewriteRulesEnabled:  utils.Convert_bool_To_Pointer_bool(false),
							UseExternalAddressForHostAndSni:     utils.Convert_bool_To_Pointer_bool(false),
							LinkedApplications:                  []string{"linked-application"},
						},
						HttpRequestCustomizationSettings: &accessv1.HttpRequestCustomizationSettings{
							HeaderCustomization: map[string]string{
								"X-Forwarded-For":   "$SOURCEIP$",
								"X-Forwarded-Host":  "$ORIGINALHOST$",
								"X-Forwarded-Proto": "$PROTOCOL$",
							},
						},
					},
					Status: accessv1.CommonApplicationStatus{
						Id: "uuid",
					},
				},
			},
			expected: &model.Application{
				ID:       "uuid",
				Type:     model.ApplicationType(model.HTTP),
				SubType:  model.ApplicationSubType(model.CustomDomain),
				ToDelete: false,
				CommonApplicationParams: model.CommonApplicationParams{
					IsVisible:             true,
					IsNotificationEnabled: true,
					Enabled:               true,
					Name:                  "my-application",
					SiteName:              "my-site",
					AccessPoliciesNames:   []string{"access-policy"},
					ActivityPoliciesNames: []string{"activity-policy"},
				},
				ConnectionSettings: &model.ConnectionSettings{
					InternalAddress:       "http://my-service-name.service-namespace:80",
					Subdomain:             "",
					CustomExternalAddress: "customExternalAddress",
					CustomRootPath:        "/root/path",
					HealthUrl:             "/health",
					HealthMethod:          "GET",
					CustomSSLCertificate:  "",
					WildcardPrivateKey:    "",
				},
				HttpLinkTranslationSettings: &model.HttpLinkTranslationSettings{
					IsDefaultContentRewriteRulesEnabled: false,
					IsDefaultHeaderRewriteRulesEnabled:  false,
					UseExternalAddressForHostAndSni:     false,
					LinkedApplications:                  []string{"linked-application"},
				},
				HttpRequestCustomizationSettings: &model.HttpRequestCustomizationSettings{
					HeaderCustomization: map[string]string{
						"X-Forwarded-For":   "$SOURCEIP$",
						"X-Forwarded-Host":  "$ORIGINALHOST$",
						"X-Forwarded-Proto": "$PROTOCOL$",
					},
				},
			},
		},
		{
			name: "default flow",
			args: args{
				application: &accessv1.HttpApplication{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-application",
					},
					Spec: accessv1.HttpApplicationSpec{
						CommonApplicationParams: accessv1.CommonApplicationParams{
							SiteName: "my-site",
						},
						Service: accessv1.Service{
							Name:      "my-service-name",
							Namespace: "service-namespace",
							Port:      "80",
						},
						HttpConnectionSettings:           nil,
						HttpLinkTranslationSettings:      nil,
						HttpRequestCustomizationSettings: nil,
					},
					Status: accessv1.CommonApplicationStatus{
						Id: "uuid",
					},
				},
			},
			expected: &model.Application{
				ID:       "uuid",
				Type:     model.ApplicationType(model.HTTP),
				SubType:  model.DefaultSubType,
				ToDelete: false,
				CommonApplicationParams: model.CommonApplicationParams{
					IsVisible:             true,
					IsNotificationEnabled: false,
					Enabled:               true,
					Name:                  "my-application",
					SiteName:              "my-site",
					AccessPoliciesNames:   nil,
					ActivityPoliciesNames: nil,
				},
				ConnectionSettings: &model.ConnectionSettings{
					InternalAddress: "http://my-service-name.service-namespace:80",
				},
				HttpLinkTranslationSettings:      nil,
				HttpRequestCustomizationSettings: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpApplicationTypeConverter{}
			got, err := a.ConvertToModel(tt.args.application)
			require.Equal(t, tt.errorOutput, err)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expected.CommonApplicationParams, got.CommonApplicationParams)
			assert.Equal(t, tt.expected.ConnectionSettings, got.ConnectionSettings)
			assert.Equal(t, tt.expected.HttpLinkTranslationSettings, got.HttpLinkTranslationSettings)
			assert.Equal(t, tt.expected.HttpRequestCustomizationSettings, got.HttpRequestCustomizationSettings)
		})
	}
}
