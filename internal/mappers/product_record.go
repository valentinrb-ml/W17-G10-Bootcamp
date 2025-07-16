package mappers

import "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/productrecord"

// ProductRecordRequestToDomain convert request to domain
func ProductRecordRequestToDomain(req productrecord.ProductRecordRequest) productrecord.ProductRecord {
	return productrecord.ProductRecord{
		ID:                0, // It is assigned in the DB
		ProductRecordCore: req.Data,
	}
}

// ProductRecordReportToResponse converts the report to response
func ProductRecordReportToResponse(reports []productrecord.ProductRecordReport) productrecord.ProductRecordsReportResponse {
	return productrecord.ProductRecordsReportResponse{
		Data: reports,
	}
}
