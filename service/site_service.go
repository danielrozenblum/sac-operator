package service

import (
	"context"

	"errors"

	"bitbucket.org/accezz-io/sac-operator/model"
)

var SiteAlreadyExist = errors.New("site already exist")

type SiteService interface {
	Reconcile(ctx context.Context, site *model.Site) error
}
