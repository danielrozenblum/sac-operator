package service

import (
	"context"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"

	"github.com/google/uuid"

	"errors"

	"bitbucket.org/accezz-io/sac-operator/model"
)

var SiteAlreadyExist = errors.New("site already exist")

type SiteService interface {
	Create(ctx context.Context, site *model.Site, siteCRD *accessv1.Site) error
	Delete(ctx context.Context, id uuid.UUID) error
}
