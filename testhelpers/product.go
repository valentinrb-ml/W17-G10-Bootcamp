package testhelpers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
)

// MustJSON serialises v to JSON or panics.
func MustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// NewRequest builds an *http.Request with JSON Content-Type.
func NewRequest(tb testing.TB, method, url string, body io.Reader) *http.Request {
	tb.Helper()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// DoRawRequest executes h over the supplied request and returns the recorder.
func DoRawRequest(tb testing.TB, req *http.Request, h http.HandlerFunc) *httptest.ResponseRecorder {
	tb.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// BuildProduct builds a valid domain Product. id==0 implies “create”.
func BuildProduct(id int) models.Product {
	seller := 200
	return models.Product{
		ID:          id,
		Code:        "CODE",
		Description: "desc",
		Dimensions: models.Dimensions{
			Width:  1,
			Height: 1,
			Length: 1,
		},
		NetWeight: 10,
		Expiration: models.Expiration{
			Rate:                    5,
			RecommendedFreezingTemp: 2,
			FreezingRate:            6,
		},
		ProductType: 100,
		SellerID:    &seller,
	}
}
