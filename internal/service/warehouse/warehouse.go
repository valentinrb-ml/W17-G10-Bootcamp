package service

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

// Create creates a new warehouse by delegating to the repository layer
// Returns the created warehouse or an error if the operation fails
func (s *WarehouseDefault) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
	return s.rp.Create(ctx, w)
}

// FindAll retrieves all warehouses from the repository and sorts them by ID
// Returns a slice of warehouses sorted by ID in ascending order or an error if the operation fails
func (s *WarehouseDefault) FindAll(ctx context.Context) ([]warehouse.Warehouse, error) {
	return s.rp.FindAll(ctx)
}

// FindById retrieves a specific warehouse by its ID from the repository
// Returns the warehouse if found or an error if not found or operation fails
func (s *WarehouseDefault) FindById(ctx context.Context, id int) (*warehouse.Warehouse, error) {
	return s.rp.FindById(ctx, id)
}

// Update modifies an existing warehouse with the provided patch data
// First retrieves the existing warehouse, validates minimum capacity if provided,
// applies the patch, and then updates the warehouse in the repository
// Returns the updated warehouse or an error if validation fails or operation fails
func (s *WarehouseDefault) Update(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, error) {
	existing, err := s.rp.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if patch.MinimumCapacity != nil {
		if err := validators.ValidateMinimumCapacity(*patch.MinimumCapacity); err != nil {
			return nil, err
		}
	}

	mappers.ApplyWarehousePatch(existing, patch)

	updated, errRepo := s.rp.Update(ctx, id, *existing)
	if errRepo != nil {
		return nil, errRepo
	}
	return updated, nil
}

// Delete removes a warehouse from the repository by its ID
// First verifies the warehouse exists, then deletes it from the repository
// Returns an error if the warehouse doesn't exist or operation fails
func (s *WarehouseDefault) Delete(ctx context.Context, id int) error {
	return s.rp.Delete(ctx, id)
}
