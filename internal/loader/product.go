package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product"
	"io"
	"os"
	"path/filepath"
)

func serviceErr(code int, internal error, overrideMsg string) error {
	e := api.ServiceErrors[code]
	if overrideMsg != "" {
		e.Message = overrideMsg
	}
	e.InternalError = internal
	return e
}

type ProductLoader interface {
	Load(ctx context.Context) (map[int]product.Product, error)
}

type JSONFileProductLoader struct {
	filePath string
}

func NewJSONFileProductLoader(path string) (*JSONFileProductLoader, error) {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, serviceErr(api.ErrNotFound, err, "products file not found")
	}
	if info.IsDir() {
		return nil, serviceErr(api.ErrBadRequest, fmt.Errorf("%s is a directory", cleanPath), "path must be a file")
	}

	return &JSONFileProductLoader{filePath: cleanPath}, nil
}

func (l *JSONFileProductLoader) Load(_ context.Context) (map[int]product.Product, error) {
	f, err := os.Open(l.filePath)
	if err != nil {
		return nil, serviceErr(api.ErrNotFound, err, "unable to open products file")
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, serviceErr(api.ErrInternalServer, err, "reading products file")
	}

	var dtoList []product.ProductResponse
	if err := json.Unmarshal(raw, &dtoList); err != nil {
		return nil, serviceErr(api.ErrInternalServer, err, "decoding products JSON")
	}

	// Convierte slice -> map y valida duplicados
	result := make(map[int]product.Product, len(dtoList))
	for _, p := range dtoList {
		if _, exists := result[p.ID]; exists {
			return nil, serviceErr(api.ErrBadRequest, nil,
				fmt.Sprintf("duplicated product id %d in JSON file", p.ID))
		}
		result[p.ID] = mappers.ResponseToDomain(p)
	}

	return result, nil
}
