package part

import (
	"context"
	"slices"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	"github.com/Medveddo/rocket-science/inventory/internal/repository/converter"
)

func (r *partsRepository) ListParts(_ context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.Part

	for _, part := range r.data {
		if filter != nil {
			if len(filter.UUIDs) > 0 && !slices.Contains(filter.UUIDs, part.UUID) {
				continue
			}
			if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
				continue
			}
			if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, model.Category(part.Category)) {
				continue
			}
			if len(filter.ManufacturerCountries) > 0 && (part.Manufacturer == nil || !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country)) {
				continue
			}
			if len(filter.Tags) > 0 {
				skip := false
				for _, tag := range filter.Tags {
					if !slices.Contains(part.Tags, tag) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
			}
		}
		part := converter.RepoPartToPart(*part)
		result = append(result, part)
	}
	return result, nil
}
