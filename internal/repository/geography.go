package repository

import (
	"context"

	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/geography"
)

const (
	queryCountryCreate   = `INSERT INTO countries (name) VALUES (?)`
	queryCountryFindById = `SELECT id, name FROM countries WHERE LOWER(name) = LOWER(?)`

	queryProvinceCreate   = `INSERT INTO provinces (name, country_id) VALUES (?, ?)`
	queryProvinceFindById = `SELECT id, name, country_id FROM provinces WHERE LOWER(name) = LOWER(?) AND country_id = ?`

	queryLocalityCreate   = `INSERT INTO localities (name, province_id) VALUES (?, ?)`
	queryLocalityFindById = `SELECT id, name, province_id FROM localities WHERE LOWER(name) = LOWER(?) AND province_id = ?`
)

func (r *geographyRepository) CreateCountry(ctx context.Context, exec Executor, c models.Country) (*models.Country, error) {
	res, err := exec.ExecContext(ctx, queryCountryCreate, c.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.Id = int(id)
	return &c, nil
}

func (r *geographyRepository) FindCountryByName(ctx context.Context, exec Executor, name string) (*models.Country, error) {
	var country models.Country
	err := exec.QueryRowContext(ctx, queryCountryFindById, name).Scan(&country.Id, &country.Name)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *geographyRepository) CreateProvince(ctx context.Context, exec Executor, p models.Province) (*models.Province, error) {
	res, err := exec.ExecContext(ctx, queryProvinceCreate, p.Name, p.CountryId)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	p.Id = int(id)
	return &p, nil
}

func (r *geographyRepository) FindProvinceByName(ctx context.Context, exec Executor, name string, countryId int) (*models.Province, error) {
	var province models.Province
	err := exec.QueryRowContext(ctx, queryProvinceFindById, name, countryId).Scan(&province.Id, &province.Name, &province.CountryId)
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (r *geographyRepository) CreateLocality(ctx context.Context, exec Executor, l models.Locality) (*models.Locality, error) {
	res, err := exec.ExecContext(ctx, queryLocalityCreate, l.Name, l.ProvinceId)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	l.Id = int(id)
	return &l, nil
}

func (r *geographyRepository) FindLocalityByName(ctx context.Context, exec Executor, name string, provinceId int) (*models.Locality, error) {
	var locality models.Locality
	err := exec.QueryRowContext(ctx, queryLocalityFindById, name, provinceId).Scan(&locality.Id, &locality.Name, &locality.ProvinceId)
	if err != nil {
		return nil, err
	}
	return &locality, nil
}
