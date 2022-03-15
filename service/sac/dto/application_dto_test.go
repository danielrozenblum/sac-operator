package dto

import (
	"testing"

	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/stretchr/testify/assert"
)

func TestConvertFromApplicationModel(t *testing.T) {
	// given
	name := "test12345"
	applicationModel := model.NewApplicationBuilder().WithName(name).Build()

	// when
	result := FromApplicationModel(applicationModel)

	// then
	assert.Equal(t, applicationModel.ID, result.ID)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, applicationModel.Type, result.Type)
	assert.Equal(t, applicationModel.SubType, result.SubType)
	assert.Equal(t, applicationModel.InternalAddress, result.ConnectionSettings.InternalAddress)
}

//func TestConvertToApplicationModel(t *testing.T) {
//	// given
//	applicationDTO := NewApplicationDTOBuilder().Build()
//	siteId := "2f45a8d6-5656-40f0-8642-d9c7bb35a076"
//
//	// when
//	result := ToApplicationModel(applicationDTO, siteId)
//
//	// then
//	assert.Equal(t, applicationDTO.ID, result.ID)
//	assert.Equal(t, applicationDTO.Name, result.Name)
//	assert.Equal(t, applicationDTO.Type, result.Type)
//	assert.Equal(t, applicationDTO.SubType, result.SubType)
//	assert.Equal(t, applicationDTO.ConnectionSettings.InternalAddress, result.InternalAddress)
//	assert.Equal(t, siteId, result.SiteName)
//}

func TestMergeApplication(t *testing.T) {
	// given
	existingApplicationDTO := NewApplicationDTOBuilder().WithIsVisible(false).Build()
	updatedApplicationDTO := NewApplicationDTOBuilder().WithName("new-name").Build()

	// when
	result := MergeApplication(existingApplicationDTO, updatedApplicationDTO)

	// then
	assert.Equal(t, existingApplicationDTO.ID, result.ID)
	assert.Equal(t, updatedApplicationDTO.Name, result.Name)
	assert.Equal(t, false, result.IsVisible)
}
