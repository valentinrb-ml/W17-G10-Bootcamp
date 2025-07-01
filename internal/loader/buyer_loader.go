package loader

import (
	"encoding/json"
	"os"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models"
)

type BuyerLoader interface {
	Load() (v map[int]models.Buyer, err error)
}

func NewBuyerJSONFile(path string) *BuyerJSONFile {
	return &BuyerJSONFile{
		path: path,
	}
}

type BuyerJSONFile struct {
	path string
}

func (l *BuyerJSONFile) Load() (b map[int]models.Buyer, err error) {
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	var buyersJSON []models.Buyer
	err = json.NewDecoder(file).Decode(&buyersJSON)
	if err != nil {
		return
	}

	b = make(map[int]models.Buyer)
	for _, buyer := range buyersJSON {
		b[buyer.Id] = models.Buyer{
			Id:           buyer.Id,
			CardNumberId: buyer.CardNumberId,
			FirstName:    buyer.FirstName,
			LastName:     buyer.LastName,
		}
	}

	return
}
