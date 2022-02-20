package converter

import (
	"testing"
	"time"

	"bitbucket.org/accezz-io/sac-operator/service"

	"github.com/stretchr/testify/assert"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
)

func TestSiteConverter_ConvertFromServiceModel(t *testing.T) {
	type args struct {
		site *service.SiteReconcileOutput
	}
	tests := []struct {
		name string
		args args
		want accessv1.SiteStatus
	}{
		{
			name: "happy-flow",
			args: args{
				site: &service.SiteReconcileOutput{
					Deleted:   false,
					SACSiteID: "51f33785-434d-41cf-8eae-7c07f43afbe1",
					HealthyConnectors: []service.Connectors{{
						CreatedTimestamp: time.Time{},
						DeploymentName:   "dep1",
						SacID:            "uuid1",
					}},
					UnHealthyConnectors: []service.Connectors{{
						CreatedTimestamp: time.Time{},
						DeploymentName:   "dep2",
						SacID:            "uuid2",
					}},
				},
			},
			want: accessv1.SiteStatus{
				ID: "51f33785-434d-41cf-8eae-7c07f43afbe1",
				HealthyConnectors: map[string]string{
					"dep1": "uuid1",
				},
				UnHealthyConnectors: map[string]string{
					"dep2": "uuid2",
				},
				NumberOfHealthyConnectors: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SiteConverter{}
			site := s.ConvertFromServiceOutput(tt.args.site)
			assert.Equal(t, tt.want, site)
		})
	}
}
