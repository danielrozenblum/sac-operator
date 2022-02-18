package service

import (
	"context"
	"errors"

	"bitbucket.org/accezz-io/sac-operator/model"
)

var SiteAlreadyExist = errors.New("site already exist")

type SiteReconcileOutput struct {
	Deleted             bool
	SACSiteID           string
	HealthConnectors    map[string]string
	UnHealthyConnectors map[string]string
}

type SiteService interface {
	Reconcile(ctx context.Context, site *model.Site) (*SiteReconcileOutput, error)
}
