package testhelpers

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/buyer"
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

func DummyRequestBuyer() models.RequestBuyer {
	cardNumber := "CARD-001"
	firstName := "John"
	lastName := "Doe"

	return models.RequestBuyer{
		CardNumberId: &cardNumber,
		FirstName:    &firstName,
		LastName:     &lastName,
	}
}
func DummyResponseBuyer() models.ResponseBuyer {
	return models.ResponseBuyer(BuyersDummyMap[1])
}

var BuyersDummyMap = map[int]models.Buyer{
	1: {
		Id:           1,
		CardNumberId: "CARD-001",
		FirstName:    "John",
		LastName:     "Doe",
	},
	2: {
		Id:           2,
		CardNumberId: "CARD-002",
		FirstName:    "Jane",
		LastName:     "Smith",
	},
	3: {
		Id:           3,
		CardNumberId: "CARD-003",
		FirstName:    "Bob",
		LastName:     "Johnson",
	},
}

func FindAllBuyersDummy() []models.Buyer {
	out := make([]models.Buyer, 0, len(BuyersDummyMap))
	for i := 1; i <= len(BuyersDummyMap); i++ {
		if buyer, ok := BuyersDummyMap[i]; ok {
			out = append(out, buyer)
		}
	}

	return out
}

func FindAllBuyersResponseDummy() []models.ResponseBuyer {
	var keys []int
	for k := range BuyersDummyMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var buyersList []models.Buyer
	for _, k := range keys {
		buyersList = append(buyersList, BuyersDummyMap[k])
	}

	mb := mappers.ToResponseBuyerList(buyersList)
	return mb
}

// Ptr helper for expected pointer struct
func PtrBuyer[T any](v T) *T { return &v }
