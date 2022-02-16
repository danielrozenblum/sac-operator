package repository

import (
	"context"
)

const siteFinalizerName = "site.access.secure-access-cloud.symantec.com/finalizer"

type Repository interface {
	UpdateNewSite(ctx context.Context, siteName, id string) error
	UpdateDeleteSite(ctx context.Context, siteName string) error
}
