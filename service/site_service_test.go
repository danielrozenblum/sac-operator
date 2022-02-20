package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_sortConnectorsBtCreatedTimestamp(t *testing.T) {
	baseTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name         string
		inConnector  []Connectors
		outConnector []Connectors
	}{
		{
			name: "",
			inConnector: []Connectors{
				{
					CreatedTimestamp: baseTime.Add(3 * time.Second),
					DeploymentName:   "third",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime.Add(1 * time.Second),
					DeploymentName:   "first",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime,
					DeploymentName:   "baseTime",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime.Add(2 * time.Second),
					DeploymentName:   "second",
					SacID:            "uuid",
				},
			},
			outConnector: []Connectors{
				{
					CreatedTimestamp: baseTime,
					DeploymentName:   "baseTime",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime.Add(1 * time.Second),
					DeploymentName:   "first",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime.Add(2 * time.Second),
					DeploymentName:   "second",
					SacID:            "uuid",
				},
				{
					CreatedTimestamp: baseTime.Add(3 * time.Second),
					DeploymentName:   "third",
					SacID:            "uuid",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortConnectorsByOldestFirst(tt.inConnector)
			assert.Equal(t, tt.inConnector, tt.outConnector)
		})
	}
}
