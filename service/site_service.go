package service

import (
	"context"

	"github.com/google/uuid"

	"errors"

	"bitbucket.org/accezz-io/sac-operator/model"
)

var SiteAlreadyExist = errors.New("site already exist")

type SiteService interface {
	Create(ctx context.Context, site *model.Site) error
	Delete(ctx context.Context, id uuid.UUID) error
	Reconcile(ctx context.Context, site *model.Site) error
}
