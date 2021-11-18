package converter

import (
	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
)

type ApplicationTypeConverter struct{}

func NewApplicationTypeConverter() *ApplicationTypeConverter {
	return &ApplicationTypeConverter{}
}

func (a *ApplicationTypeConverter) ConvertToModel(application accessv1.Application) *model.Application {
	return nil
}

func (a *ApplicationTypeConverter) ConvertFromModel(application model.Application) *accessv1.Application {
	return nil
}
