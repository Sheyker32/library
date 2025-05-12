package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func CheckIfTableHasRecords(db *sqlx.DB, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s LIMIT 1)", tableName)
	err := db.Get(&exists, query)
	return exists, err
}
