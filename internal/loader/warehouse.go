package loader

import (
	"encoding/json"
	"os"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseLoader interface {
	Load() (v map[int]models.Warehouse, err error)
}

// NewWarehouseJSONFile is a function that returns a new instance of WarehouseJSONFile
func NewWarehouseJSONFile(path string) *WarehouseJSONFile {
	return &WarehouseJSONFile{
		path: path,
	}
}

// WarehouseJSONFile is a struct that implements the LoaderWarehouse interface
type WarehouseJSONFile struct {
	// path is the path to the file that contains the Warehouse in JSON format
	path string
}

// Load is a method that loads the Warehouse
func (l *WarehouseJSONFile) Load() (s map[int]models.Warehouse, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var WarehouseJSON []models.Warehouse
	err = json.NewDecoder(file).Decode(&WarehouseJSON)
	if err != nil {
		return
	}

	// serialize Warehouse
	s = make(map[int]models.Warehouse)
	for _, se := range WarehouseJSON {
		s[se.Id] = models.Warehouse{
			Id:          se.Id,
			Address:         se.Address,
			WarehouseCode: se.WarehouseCode,
			MinimumTemperature:     se.MinimumTemperature,
			MinimumCapacity:   se.MinimumCapacity,
			Telephone:         se.Telephone,
		}
	}

	return
}