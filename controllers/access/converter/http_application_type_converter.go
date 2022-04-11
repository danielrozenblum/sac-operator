package converter

import (
	"fmt"

	"bitbucket.org/accezz-io/sac-operator/utils"

	"github.com/jinzhu/copier"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CommonParamsConverter struct{}

func (c *CommonParamsConverter) copyCommonParams(params accessv1.CommonApplicationParams, applicationParams *model.CommonApplicationParams) error {
	applicationParams.IsVisible = utils.Convert_Pointer_bool_To_bool_with_default(params.IsVisible, true)
	applicationParams.Enabled = utils.Convert_Pointer_bool_To_bool_with_default(params.Enabled, true)
	applicationParams.IsNotificationEnabled = utils.Convert_Pointer_bool_To_bool_with_default(params.IsNotificationEnabled, false)
	applicationParams.SiteName = params.SiteName
	applicationParams.AccessPoliciesNames = params.AccessPoliciesNames
	applicationParams.ActivityPoliciesNames = params.ActivityPoliciesNames

	return nil
}

func convertToSchema(applicationType model.ApplicationType, service accessv1.Service) string {
	if service.Schema != "" {
		return service.Schema
	}

	switch applicationType {
	case model.SSH, model.DynamicSSH, model.RDP, model.TCP:
		return "tcp"
	case model.HTTP:
		{
			switch service.Port {
			case "443", "8443":
				return "https"
			default:
				return "http"
			}
		}
	default:
		return "http"
	}
}

type HttpApplicationTypeConverter struct {
	*CommonParamsConverter
}

func NewHttpApplicationTypeConverter() *HttpApplicationTypeConverter {
	return &HttpApplicationTypeConverter{
		&CommonParamsConverter{},
	}
}

func (a *HttpApplicationTypeConverter) Validate(application *accessv1.HttpApplication) error {

	if application.Spec.Service.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if application.Spec.Service.Port == "" {
		return fmt.Errorf("service port cannot be empty")
	}

	return nil
}

func (a *HttpApplicationTypeConverter) ConvertToModel(application *accessv1.HttpApplication) (*model.Application, error) {

	output := &model.Application{
		ID:       application.Status.Id,
		Type:     model.HTTP,
		SubType:  utils.GetApplicationSubTypeOrDefault(application.Spec.SubType, model.DefaultSubType),
		ToDelete: !application.ObjectMeta.DeletionTimestamp.IsZero(),
		ConnectionSettings: &model.ConnectionSettings{
			InternalAddress: a.convertToInternalAddress(application.Spec.Service, application.Namespace),
		},
	}

	var err error

	if application.Spec.HttpConnectionSettings != nil {
		err = copier.Copy(&output.ConnectionSettings, application.Spec.HttpConnectionSettings)
		if err != nil {
			return output, err
		}
	}

	if application.Spec.HttpLinkTranslationSettings != nil {
		httpLinkTranslationSettings := &model.HttpLinkTranslationSettings{}
		err = copier.Copy(httpLinkTranslationSettings, application.Spec.HttpLinkTranslationSettings)
		if err != nil {
			return output, err
		}
		output.HttpLinkTranslationSettings = httpLinkTranslationSettings
	}

	if application.Spec.HttpRequestCustomizationSettings != nil {
		httpRequestCustomizationSettings := &model.HttpRequestCustomizationSettings{}
		err = copier.Copy(httpRequestCustomizationSettings, application.Spec.HttpRequestCustomizationSettings)
		if err != nil {
			return output, err
		}
		output.HttpRequestCustomizationSettings = httpRequestCustomizationSettings
	}
	commonParams := model.CommonApplicationParams{
		Name: application.Name,
	}
	err = a.copyCommonParams(application.Spec.CommonApplicationParams, &commonParams)
	if err != nil {
		return output, err
	}

	output.CommonApplicationParams = commonParams

	return output, nil
}

func (a *HttpApplicationTypeConverter) convertToInternalAddress(service accessv1.Service, applicationNamespace string) string {

	schema := convertToSchema(model.HTTP, service)
	namespace := applicationNamespace
	if service.Namespace != "" {
		namespace = service.Namespace
	}
	if service.Port == "" {
		return fmt.Sprintf("%s://%s.%s", schema, service.Name, namespace)
	}
	return fmt.Sprintf("%s://%s.%s:%s", schema, service.Name, namespace, service.Port)
}

func (a HttpApplicationTypeConverter) ConvertFromServiceOutput(output *service.ApplicationReconcileOutput) accessv1.CommonApplicationStatus {
	return accessv1.CommonApplicationStatus{
		Id:         output.SACApplicationID,
		ModifiedOn: metav1.Now(),
	}
}
