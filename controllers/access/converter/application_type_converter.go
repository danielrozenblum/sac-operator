package converter

import (
	"fmt"
	"strings"

	"bitbucket.org/accezz-io/sac-operator/service"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/utils"
)

type ApplicationTypeConverter struct{}

func NewApplicationTypeConverter() *ApplicationTypeConverter {
	return &ApplicationTypeConverter{}
}

func (a *ApplicationTypeConverter) ConvertToModel(application *accessv1.Application) *model.Application {

	applicationName := utils.GetStringPtrValueOrDefault(application.Spec.Name, application.Namespace+"-"+application.Spec.Service.Name)
	applicationType := utils.GetApplicationTypeOrDefault(application.Spec.Type, model.DefaultType)
	applicationSubType := utils.GetApplicationSubTypeOrDefault(application.Spec.SubType, model.DefaultSubType)

	return &model.Application{
		ID:               application.Status.Id,
		Name:             a.convertToValidSACApplicationName(applicationName),
		Type:             applicationType,
		SubType:          applicationSubType,
		InternalAddress:  a.convertToInternalAddress(applicationType, application.Spec.Service),
		Site:             application.Spec.Site,
		AccessPolicies:   application.Spec.AccessPolicies,
		ActivityPolicies: application.Spec.ActivityPolicies,
		ToDelete:         !application.ObjectMeta.DeletionTimestamp.IsZero(),
	}
}

func (a *ApplicationTypeConverter) convertToValidSACApplicationName(value string) string {
	result := strings.ReplaceAll(value, " ", "-")

	if len(result) > 64 {
		return result[0:63]
	}

	return result
}

func (a *ApplicationTypeConverter) convertToInternalAddress(applicationType model.ApplicationType, service accessv1.Service) string {
	schema := a.convertToSchema(applicationType, service)
	namespace := "default"
	if service.Namespace != "" {
		namespace = service.Namespace
	}
	return fmt.Sprintf("%s://%s.%s:%s", schema, service.Name, namespace, service.Port)
}

func (a *ApplicationTypeConverter) convertToSchema(applicationType model.ApplicationType, service accessv1.Service) string {
	if service.Schema != nil {
		return *service.Schema
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

func (a ApplicationTypeConverter) ConvertFromServiceOutput(output *service.ApplicationReconcileOutput) accessv1.ApplicationStatus {
	return accessv1.ApplicationStatus{}
}
