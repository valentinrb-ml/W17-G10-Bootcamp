package mappers

import productrecord "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/product_record"

// ProductRecordRequestToDomain convert request to domain
func ProductRecordRequestToDomain(req productrecord.ProductRecordRequest) productrecord.ProductRecord {
	return productrecord.ProductRecord{
		ID:                0, // It is assigned in the DB
		ProductRecordCore: req.Data,
	}
}
