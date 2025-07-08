package part

import (
	"context"

	model "github.com/Medveddo/rocket-science/inventory/internal/model"
)

func (r *partService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := r.repository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return part, err
}
