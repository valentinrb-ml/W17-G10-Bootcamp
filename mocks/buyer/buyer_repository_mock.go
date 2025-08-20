package mocks

import (
	"context"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

type BuyerRepositoryMocks struct {
	MockCreate           func(ctx context.Context, b models.Buyer) (*models.Buyer, error)
	MockFindAll          func(ctx context.Context) ([]models.Buyer, error)
	MockFindByID         func(ctx context.Context, id int) (*models.Buyer, error)
	MockFindById         func(ctx context.Context, id int) (*models.Buyer, error)
	MockFindByCardNumber func(ctx context.Context, cardNumber string) (*models.Buyer, error)
	MockUpdate           func(ctx context.Context, id int, b models.Buyer) error
	MockDelete           func(ctx context.Context, id int) error
	MockCardNumberExists func(ctx context.Context, cardNumber string, id int) bool
	MockSetLogger        func(l logger.Logger)
}

func (m *BuyerRepositoryMocks) Create(ctx context.Context, b models.Buyer) (*models.Buyer, error) {
	return m.MockCreate(ctx, b)
}

func (m *BuyerRepositoryMocks) FindAll(ctx context.Context) ([]models.Buyer, error) {
	return m.MockFindAll(ctx)
}

func (m *BuyerRepositoryMocks) FindByID(ctx context.Context, id int) (*models.Buyer, error) {
	return m.MockFindByID(ctx, id)
}

func (m *BuyerRepositoryMocks) FindById(ctx context.Context, id int) (*models.Buyer, error) {
	return m.MockFindById(ctx, id)
}

func (m *BuyerRepositoryMocks) FindByCardNumber(ctx context.Context, cardNumber string) (*models.Buyer, error) {
	return m.MockFindByCardNumber(ctx, cardNumber)
}

func (m *BuyerRepositoryMocks) Update(ctx context.Context, id int, b models.Buyer) error {
	return m.MockUpdate(ctx, id, b)
}

func (m *BuyerRepositoryMocks) Delete(ctx context.Context, id int) error {
	return m.MockDelete(ctx, id)
}

func (m *BuyerRepositoryMocks) CardNumberExists(ctx context.Context, cardNumber string, id int) bool {
	return m.MockCardNumberExists(ctx, cardNumber, id)
}

func (m *BuyerRepositoryMocks) SetLogger(l logger.Logger) {
	if m.MockSetLogger != nil {
		m.MockSetLogger(l)
	}
}
