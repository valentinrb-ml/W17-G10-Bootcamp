package mocks

import (
    "context"

    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/warehouse"
)

type WarehouseServiceMock struct {
    // Funciones mock
    FuncCreate   func(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error)
    FuncFindAll  func(ctx context.Context) ([]warehouse.Warehouse, error)
    FuncFindById func(ctx context.Context, id int) (*warehouse.Warehouse, error)
    FuncUpdate   func(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, error)
    FuncDelete   func(ctx context.Context, id int) error

    // Contadores para verificar llamadas
    CreateCallCount   int
    FindAllCallCount  int
    FindByIdCallCount int
    UpdateCallCount   int
    DeleteCallCount   int

    // Parámetros recibidos (para verificar que se llamó con los valores correctos)
    CreateCalls   []CreateCall
    FindByIdCalls []FindByIdCall
    UpdateCalls   []UpdateCall
    DeleteCalls   []DeleteCall
}

// Estructuras para capturar parámetros de las llamadas
type CreateCall struct {
    Ctx       context.Context
    Warehouse warehouse.Warehouse
}

type FindByIdCall struct {
    Ctx context.Context
    Id  int
}

type UpdateCall struct {
    Ctx   context.Context
    Id    int
    Patch warehouse.WarehousePatchDTO
}

type DeleteCall struct {
    Ctx context.Context
    Id  int
}

// Implementación de la interfaz WarehouseServiceInterface
func (m *WarehouseServiceMock) Create(ctx context.Context, w warehouse.Warehouse) (*warehouse.Warehouse, error) {
    m.CreateCallCount++
    m.CreateCalls = append(m.CreateCalls, CreateCall{Ctx: ctx, Warehouse: w})
    
    if m.FuncCreate != nil {
        return m.FuncCreate(ctx, w)
    }
    return nil, nil
}

func (m *WarehouseServiceMock) FindAll(ctx context.Context) ([]warehouse.Warehouse, error) {
    m.FindAllCallCount++
    
    if m.FuncFindAll != nil {
        return m.FuncFindAll(ctx)
    }
    return nil, nil
}

func (m *WarehouseServiceMock) FindById(ctx context.Context, id int) (*warehouse.Warehouse, error) {
    m.FindByIdCallCount++
    m.FindByIdCalls = append(m.FindByIdCalls, FindByIdCall{Ctx: ctx, Id: id})
    
    if m.FuncFindById != nil {
        return m.FuncFindById(ctx, id)
    }
    return nil, nil
}

func (m *WarehouseServiceMock) Update(ctx context.Context, id int, patch warehouse.WarehousePatchDTO) (*warehouse.Warehouse, error) {
    m.UpdateCallCount++
    m.UpdateCalls = append(m.UpdateCalls, UpdateCall{Ctx: ctx, Id: id, Patch: patch})
    
    if m.FuncUpdate != nil {
        return m.FuncUpdate(ctx, id, patch)
    }
    return nil, nil
}

func (m *WarehouseServiceMock) Delete(ctx context.Context, id int) error {
    m.DeleteCallCount++
    m.DeleteCalls = append(m.DeleteCalls, DeleteCall{Ctx: ctx, Id: id})
    
    if m.FuncDelete != nil {
        return m.FuncDelete(ctx, id)
    }
    return nil
}

// Métodos helper para verificaciones en tests
func (m *WarehouseServiceMock) AssertCreateCalledWith(t interface{}, expectedCtx context.Context, expectedWarehouse warehouse.Warehouse) {
    // Implementar verificación usando testify/require
}

func (m *WarehouseServiceMock) AssertFindByIdCalledWith(t interface{}, expectedCtx context.Context, expectedId int) {
    // Implementar verificación usando testify/require
}

// Reset para limpiar entre tests
func (m *WarehouseServiceMock) Reset() {
    m.CreateCallCount = 0
    m.FindAllCallCount = 0
    m.FindByIdCallCount = 0
    m.UpdateCallCount = 0
    m.DeleteCallCount = 0
    
    m.CreateCalls = nil
    m.FindByIdCalls = nil
    m.UpdateCalls = nil
    m.DeleteCalls = nil
}