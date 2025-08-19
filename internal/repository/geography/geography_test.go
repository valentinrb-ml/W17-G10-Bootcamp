package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	geography "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/geography"
)

func TestGeographyRepository_BeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := geography.NewGeographyRepository(db)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	_ = tx.Rollback() // aseguramos cleanup de la tx
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGeographyRepository_CommitTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := geography.NewGeographyRepository(db)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	mock.ExpectCommit()

	err = repo.CommitTx(tx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGeographyRepository_RollbackTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := geography.NewGeographyRepository(db)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)
	mock.ExpectRollback()

	err = repo.RollbackTx(tx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGeographyRepository_GetDB(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := geography.NewGeographyRepository(db)

	dbReturned := repo.GetDB()
	assert.Equal(t, db, dbReturned)
}
