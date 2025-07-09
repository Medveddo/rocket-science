package part

import (
	"sync"
	"time"

	repo "github.com/Medveddo/rocket-science/inventory/internal/repository"
	repoModel "github.com/Medveddo/rocket-science/inventory/internal/repository/model"
	u "github.com/Medveddo/rocket-science/inventory/internal/utils"
)

var _ repo.PartRepository = (*partsRepository)(nil)

type partsRepository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Part
}

func NewRepository() *partsRepository {
	now := time.Now()

	part1 := &repoModel.Part{
		UUID:          "111e4567-e89b-12d3-a456-426614174001",
		Name:          "Hyperdrive Engine",
		Description:   "A class-9 hyperdrive engine capable of faster-than-light travel.",
		Price:         450000.00,
		StockQuantity: 3,
		Category:      repoModel.Category_CATEGORY_ENGINE,
		Dimensions: &repoModel.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 100.0,
			Weight: 500.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Hyperdrive Corp",
			Country: "USA",
			Website: "https://hyperdrive.example.com",
		},
		Tags: []string{"engine", "hyperdrive", "space"},
		Metadata: map[string]*repoModel.Value{
			"power_output":    {DoubleValue: u.ToPtr(9.5)},
			"is_experimental": {BoolValue: u.ToPtr(true)},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	part2 := &repoModel.Part{
		UUID:          "222e4567-e89b-12d3-a456-426614174002",
		Name:          "Quantum Shield Generator",
		Description:   "Advanced shield generator providing protection against cosmic radiation.",
		Price:         175000.00,
		StockQuantity: 5,
		Category:      repoModel.Category_CATEGORY_SHIELD,
		Dimensions: &repoModel.Dimensions{
			Length: 60.0,
			Width:  40.0,
			Height: 50.0,
			Weight: 150.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Quantum Tech",
			Country: "Germany",
			Website: "https://quantumtech.example.com",
		},
		Tags: []string{"shield", "quantum", "defense"},
		Metadata: map[string]*repoModel.Value{
			"energy_consumption": {DoubleValue: u.ToPtr(3.2)},
			"warranty_years":     {Int64Value: u.ToPtr(int64(5))},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	parts := make(map[string]*repoModel.Part, 2)

	parts[part1.UUID] = part1
	parts[part2.UUID] = part2

	return &partsRepository{
		data: parts,
	}
}
