package part

import (
	"context"

	model "github.com/Medveddo/rocket-science/inventory/internal/model"
)

func (r *partService) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	parts, err := r.repository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
