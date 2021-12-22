package converter

import (
	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/utils"
	"bitbucket.org/accezz-io/sac-operator/utils/typederror"
	"strings"
)

type ApplicationTypeConverter struct{}

func NewApplicationTypeConverter() *ApplicationTypeConverter {
	return &ApplicationTypeConverter{}
}

func (a *ApplicationTypeConverter) ConvertToModel(application accessv1.Application) (*model.Application, error) {

	applicationId, err := utils.FromUIDType(application.Status.Id)
	if err != nil {
		return nil, typederror.WrapError(typederror.UnrecoverableError, err)
	}

	applicationName := utils.GetStringPtrValueOrDefault(application.Spec.Name, application.Namespace+"-"+application.Spec.Service.Name)
	applicationType := utils.GetApplicationTypeOrDefault(application.Spec.Type, model.DefaultType)
	applicationSubType := utils.GetApplicationSubTypeOrDefault(application.Spec.SubType, model.DefaultSubType)

	return &model.Application{
		ID:               applicationId,
		Name:             a.convertToValidSACApplicationName(applicationName),
		Type:             applicationType,
		SubType:          applicationSubType,
		InternalAddress:  a.convertToInternalAddress(applicationType, application.Spec.Service),
		Site:             application.Spec.Site,
		AccessPolicies:   application.Spec.AccessPolicies,
		ActivityPolicies: application.Spec.ActivityPolicies,
	}, nil
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
	return schema + service.Name + ":" + service.Port
}

func (a *ApplicationTypeConverter) convertToSchema(applicationType model.ApplicationType, service accessv1.Service) string {
	if service.Schema != nil {
		return *service.Schema
	}

	switch applicationType {
	case model.SSH, model.DynamicSSH, model.RDP, model.TCP:
		return "tcp://"
	case model.HTTP:
		{
			switch service.Port {
			case "443", "8443":
				return "https://"
			default:
				return "http://"
			}
		}
	default:
		return "http://"
	}
}
