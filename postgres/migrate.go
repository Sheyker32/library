package postgres

import (
	"library/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigration(c *config.DBConfig) *migrate.Migrate {
	m, err := migrate.New(
		"file://migrations",
		"postgres://"+c.User+":"+c.Password+"@"+c.Host+":"+c.Port+"/"+c.DBName+"?sslmode="+c.SSLMode+"")

	if err != nil {
		log.Fatalf("failed to migrate db: %s", err.Error())
	}

	return m
}
