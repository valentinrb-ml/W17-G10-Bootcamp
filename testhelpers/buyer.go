package testhelpers

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
)

// BuyerBuilder permite construir buyers con valores personalizables
type BuyerBuilder struct {
	buyer models.Buyer
}

func NewBuyerBuilder() *BuyerBuilder {
	return &BuyerBuilder{
		buyer: models.Buyer{
			CardNumberId: "CARD-001",
			FirstName:    "John",
			LastName:     "Doe",
		},
	}
}

func (b *BuyerBuilder) WithID(id int) *BuyerBuilder {
	b.buyer.Id = id
	return b
}

func (b *BuyerBuilder) WithCardNumber(cardNumber string) *BuyerBuilder {
	b.buyer.CardNumberId = cardNumber
	return b
}

func (b *BuyerBuilder) WithFirstName(firstName string) *BuyerBuilder {
	b.buyer.FirstName = firstName
	return b
}

func (b *BuyerBuilder) WithLastName(lastName string) *BuyerBuilder {
	b.buyer.LastName = lastName
	return b
}

func (b *BuyerBuilder) Build() models.Buyer {
	return b.buyer
}

// CreateTestBuyer crea un buyer básico para testing
func CreateTestBuyer() models.Buyer {
	return NewBuyerBuilder().Build()
}

// CreateTestBuyerWithID crea un buyer con ID para resultados esperados
func CreateTestBuyerWithID(id int) *models.Buyer {
	b := NewBuyerBuilder().WithID(id).Build()
	return &b
}

// CreateTestBuyers crea múltiples buyers para testing
func CreateTestBuyers(count int) []models.Buyer {
	var buyers []models.Buyer
	for i := 1; i <= count; i++ {
		buyer := NewBuyerBuilder().
			WithID(i).
			WithCardNumber(fmt.Sprintf("CARD-%03d", i)).
			WithFirstName(fmt.Sprintf("FirstName%d", i)).
			WithLastName(fmt.Sprintf("LastName%d", i)).
			Build()
		buyers = append(buyers, buyer)
	}
	return buyers
}

// CreateMockDB crea una base de datos mock para testing
func CreateMockBuyerDB() (sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return mock, db
}

// CreateBuyerRows crea filas mock para resultados de consultas SQL
func CreateBuyerRows(mock sqlmock.Sqlmock, buyers []models.Buyer) *sqlmock.Rows {
	rows := mock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	for _, b := range buyers {
		rows.AddRow(b.Id, b.CardNumberId, b.FirstName, b.LastName)
	}
	return rows
}
