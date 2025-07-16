package loader

import (
	"encoding/json"
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"os"
)

type SectionLoader interface {
	Load() (v map[int]models.Section, err error)
}

type SectionJSONFile struct {
	path string
}

func NewSectionJSONFile(path string) *SectionJSONFile {
	return &SectionJSONFile{path: path}
}

// Load loads the section data from the json file.
func (l *SectionJSONFile) Load() (map[int]models.Section, error) {
	file, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sectionDoc []models.Section
	err = json.NewDecoder(file).Decode(&sectionDoc)
	if err != nil {
		return nil, fmt.Errorf("error loading section.json: %v", err)
	}

	sectionMap := make(map[int]models.Section, len(sectionDoc))

	for _, s := range sectionDoc {
		sec := models.Section{
			Id:                 s.Id,
			SectionNumber:      s.SectionNumber,
			CurrentTemperature: s.CurrentTemperature,
			MinimumTemperature: s.MinimumTemperature,
			CurrentCapacity:    s.CurrentCapacity,
			MinimumCapacity:    s.MinimumCapacity,
			MaximumCapacity:    s.MaximumCapacity,
			WarehouseId:        s.WarehouseId,
			ProductTypeId:      s.ProductTypeId,
		}
		sectionMap[s.Id] = sec

	}

	return sectionMap, nil
}
