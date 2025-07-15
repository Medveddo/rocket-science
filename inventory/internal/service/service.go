package service

import (
	"context"

	model "github.com/Medveddo/rocket-science/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error)
}
