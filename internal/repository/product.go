package repository

import (
	"context"
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]product.Product, error)
	GetByID(ctx context.Context, id int) (product.Product, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Save(ctx context.Context, p product.Product) (product.Product, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.Product, error)
}

type productMemoryRepository struct {
	db     map[int]product.Product
	nextID int
}

func NewProductRepository(db map[int]product.Product) ProductRepository {
	var highestID int
	for id := range db {
		if id > highestID {
			highestID = id
		}
	}

	return &productMemoryRepository{
		db:     db,
		nextID: highestID + 1,
	}
}

func (r *productMemoryRepository) GetAll(_ context.Context) ([]product.Product, error) {
	products := make([]product.Product, 0, len(r.db))
	for _, currentProduct := range r.db {
		products = append(products, currentProduct)
	}
	return products, nil
}

func (r *productMemoryRepository) GetByID(_ context.Context, id int) (product.Product, error) {
	currentProduct, found := r.db[id]
	if !found {
		return product.Product{}, apperrors.NotFound(fmt.Sprintf("product with id %d not found", id))
	}
	return currentProduct, nil
}

func (r *productMemoryRepository) ExistsByCode(_ context.Context, code string) (bool, error) {
	for _, currentProduct := range r.db {
		if currentProduct.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *productMemoryRepository) Save(ctx context.Context, currentProduct product.Product) (product.Product, error) {
	if currentProduct.ID == 0 { // Create
		exists, err := r.ExistsByCode(ctx, currentProduct.Code)
		if err != nil {
			return product.Product{}, apperrors.Wrap(err, "failed to check product code existence")
		}
		if exists {
			return product.Product{}, apperrors.Conflict("product code already exists")
		}

		currentProduct.ID = r.nextID
		r.nextID += 1
	} else { // Update
		_, err := r.GetByID(ctx, currentProduct.ID)
		if err != nil {
			return product.Product{}, err
		}

		for id, existing := range r.db {
			if existing.Code == currentProduct.Code && id != currentProduct.ID {
				return product.Product{}, apperrors.Conflict("product code already exists")
			}
		}
	}

	r.db[currentProduct.ID] = currentProduct
	return currentProduct, nil
}

func (r *productMemoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	delete(r.db, id)
	return nil
}

func (r *productMemoryRepository) Patch(ctx context.Context, id int, req product.ProductPatchRequest) (product.Product, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	updated := current
	if req.ProductCode != nil {
		updated.Code = *req.ProductCode
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	if req.Width != nil {
		updated.Dimensions.Width = *req.Width
	}
	if req.Height != nil {
		updated.Dimensions.Height = *req.Height
	}
	if req.Length != nil {
		updated.Dimensions.Length = *req.Length
	}
	if req.NetWeight != nil {
		updated.NetWeight = *req.NetWeight
	}
	if req.ExpirationRate != nil {
		updated.Expiration.Rate = *req.ExpirationRate
	}
	if req.RecommendedFreezingTemperature != nil {
		updated.Expiration.RecommendedFreezingTemp = *req.RecommendedFreezingTemperature
	}
	if req.FreezingRate != nil {
		updated.Expiration.FreezingRate = *req.FreezingRate
	}
	if req.ProductTypeID != nil {
		updated.ProductType = *req.ProductTypeID
	}
	if req.SellerID != nil {
		updated.SellerID = req.SellerID
	}

	if req.ProductCode != nil && *req.ProductCode != current.Code {
		exists, err := r.ExistsByCode(ctx, *req.ProductCode)
		if err != nil {
			return product.Product{}, apperrors.Wrap(err, "failed to check product code existence")
		}
		if exists {
			return product.Product{}, apperrors.Conflict("product code already exists")
		}
	}

	r.db[id] = updated
	return updated, nil
}
