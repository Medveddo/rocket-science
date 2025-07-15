package part

import (
	"github.com/Medveddo/rocket-science/inventory/internal/repository"
	service "github.com/Medveddo/rocket-science/inventory/internal/service"
)

var _ service.PartService = (*partService)(nil)

type partService struct {
	repository repository.PartRepository
}

func NewPartService(repository repository.PartRepository) *partService {
	return &partService{
		repository: repository,
	}
}
