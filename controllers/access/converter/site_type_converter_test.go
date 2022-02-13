package converter

import (
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
)

func TestSiteConverter_ConvertFromServiceModel(t *testing.T) {
	type args struct {
		site *model.Site
	}
	tests := []struct {
		name string
		args args
		want *accessv1.SiteStatus
	}{
		{
			name: "happy-flow",
			args: args{
				site: &model.Site{
					SACSiteID: getRandomUUIDPointer(),
					Connectors: []model.Connector{
						{
							ConnectorID:           getRandomUUIDPointer(),
							ConnectorDeploymentID: getRandomUUIDPointer(),
							Name:                  "test",
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SiteConverter{}
			site := s.ConvertFromServiceModel(tt.args.site)
			assert.Equal(t, string(*site.ID), tt.args.site.SACSiteID.String())
			for i := range tt.args.site.Connectors {
				assert.Equal(t, tt.args.site.Connectors[i].ConnectorID.String(), string(*site.Connectors[i].ConnectorID))
				assert.Equal(t, tt.args.site.Connectors[i].ConnectorDeploymentID.String(), string(*site.Connectors[i].PodID))
			}
		})
	}
}

func getRandomUUIDPointer() *uuid.UUID {
	id := uuid.New()
	return &id
}
