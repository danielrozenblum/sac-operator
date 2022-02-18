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
					SACSiteID: "51f33785-434d-41cf-8eae-7c07f43afbe1",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SiteConverter{}
			site := s.ConvertFromServiceOutput(tt.args.site)
			assert.Equal(t, tt.args.site.SACSiteID, site.ID)
		})
	}
}

func getRandomUUIDPointer() *uuid.UUID {
	id := uuid.New()
	return &id
}
