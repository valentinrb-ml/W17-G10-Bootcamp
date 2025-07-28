package testhelpers

import (
	"context"
	"github.com/go-chi/chi/v5"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"net/http"
)

func DummySection(id int) models.Section {
	return models.Section{
		Id:                 id,
		SectionNumber:      id*10 + 1,
		CurrentCapacity:    20,
		CurrentTemperature: 5,
		MaximumCapacity:    100,
		MinimumCapacity:    10,
		MinimumTemperature: 5,
		ProductTypeId:      2,
		WarehouseId:        1,
	}
}
func DummySectionPatch(id int) models.PatchSection {
	return models.PatchSection{
		SectionNumber:      IntPtr(id*10 + 1),
		CurrentCapacity:    IntPtr(20),
		CurrentTemperature: Float64Ptr(5),
		MaximumCapacity:    IntPtr(100),
		MinimumCapacity:    IntPtr(10),
		MinimumTemperature: Float64Ptr(5),
		ProductTypeId:      IntPtr(2),
		WarehouseId:        IntPtr(1),
	}
}

func DummySectionPost(id int) models.PostSection {
	return models.PostSection{
		SectionNumber:      id*10 + 1,
		CurrentCapacity:    20,
		CurrentTemperature: Float64Ptr(5),
		MaximumCapacity:    100,
		MinimumCapacity:    10,
		MinimumTemperature: Float64Ptr(5),
		ProductTypeId:      2,
		WarehouseId:        1,
	}
}

func DummyResponseSection(id int) models.ResponseSection {
	return models.ResponseSection{
		Id:                 id,
		SectionNumber:      id*10 + 1,
		CurrentCapacity:    20,
		CurrentTemperature: 5,
		MaximumCapacity:    100,
		MinimumCapacity:    10,
		MinimumTemperature: 5,
		ProductTypeId:      2,
		WarehouseId:        1,
	}
}

func SetChiURLParam(req *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
