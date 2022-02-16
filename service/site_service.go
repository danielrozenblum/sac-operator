package service

import (
	"context"
	"time"

	"errors"

	"bitbucket.org/accezz-io/sac-operator/model"
)

var SiteAlreadyExist = errors.New("site already exist")

type ReconcileOutput struct {
	RequeueAfter time.Duration
}

type SiteService interface {
	Reconcile(ctx context.Context, site *model.Site) (*ReconcileOutput, error)
}
