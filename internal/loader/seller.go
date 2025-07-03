package loader

import (
	"encoding/json"
	"os"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/seller"
)

type SellerLoader interface {
	Load() (v map[int]models.Seller, err error)
}

// NewSellerJSONFile is a function that returns a new instance of SellerJSONFile
func NewSellerJSONFile(path string) *SellerJSONFile {
	return &SellerJSONFile{
		path: path,
	}
}

// SellerJSONFile is a struct that implements the LoaderSeller interface
type SellerJSONFile struct {
	// path is the path to the file that contains the Sellers in JSON format
	path string
}

// Load is a method that loads the Sellers
func (l *SellerJSONFile) Load() (s map[int]models.Seller, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var SellersJSON []models.Seller
	err = json.NewDecoder(file).Decode(&SellersJSON)
	if err != nil {
		return
	}

	// serialize Sellers
	s = make(map[int]models.Seller)
	for _, se := range SellersJSON {
		s[se.Id] = models.Seller{
			Id:          se.Id,
			Cid:         se.Cid,
			CompanyName: se.CompanyName,
			Address:     se.Address,
			Telephone:   se.Telephone,
		}
	}

	return
}
