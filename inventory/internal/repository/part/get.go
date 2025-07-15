package part

import (
	"context"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	"github.com/Medveddo/rocket-science/inventory/internal/repository/converter"
)

func (r *partsRepository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartDoesNotExist
	}
	return converter.RepoPartToPart(*part), nil
}
