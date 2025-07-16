package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitMysqlDatabase() (*sql.DB, error) {
	strConn := os.Getenv("MYSQL_CONN")
	if strConn == "" {
		return nil, fmt.Errorf("MYSQL_CONN environment not set")
	}

	db, err := sql.Open("mysql", strConn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
