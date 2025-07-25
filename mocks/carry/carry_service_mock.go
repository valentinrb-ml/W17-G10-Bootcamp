package mocks

import (
    "context"
    
    "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/carry"
)

type CarryServiceMock struct {
    FuncCreate           func(ctx context.Context, c carry.Carry) (*carry.Carry, error)
    FuncGetCarriesReport func(ctx context.Context, localityID *string) (interface{}, error)
    
    // Para verificar llamadas
    CreateCallCount           int
    CreateCalls              []CreateCall
    GetCarriesReportCallCount int
    GetCarriesReportCalls    []GetCarriesReportCall
}

type CreateCall struct {
    Ctx   context.Context
    Carry carry.Carry
}

type GetCarriesReportCall struct {
    Ctx        context.Context
    LocalityID *string
}

func (m *CarryServiceMock) Create(ctx context.Context, c carry.Carry) (*carry.Carry, error) {
    m.CreateCallCount++
    m.CreateCalls = append(m.CreateCalls, CreateCall{Ctx: ctx, Carry: c})
    return m.FuncCreate(ctx, c)
}

func (m *CarryServiceMock) GetCarriesReport(ctx context.Context, localityID *string) (interface{}, error) {
    m.GetCarriesReportCallCount++
    m.GetCarriesReportCalls = append(m.GetCarriesReportCalls, GetCarriesReportCall{Ctx: ctx, LocalityID: localityID})
    return m.FuncGetCarriesReport(ctx, localityID)
}