package service

import (
	"context"
	"sort"
	"time"

	"bitbucket.org/accezz-io/sac-operator/model"
)

type Connectors struct {
	CreatedTimestamp time.Time
	DeploymentName   string
	SacID            string
}

func sortConnectorsByOldestFirst(c []Connectors) {

	sort.Slice(c, func(i, j int) bool {
		return c[i].CreatedTimestamp.Before(c[j].CreatedTimestamp)
	})

}

type SiteReconcileOutput struct {
	Deleted             bool
	SACSiteID           string
	HealthyConnectors   []Connectors
	UnHealthyConnectors []Connectors
}

type SiteService interface {
	Reconcile(ctx context.Context, site *model.Site) (*SiteReconcileOutput, error)
}
