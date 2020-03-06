package common

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // load db driver
)

var db *gorm.DB

// OpenDB opens a postgres database specified by environment variable.
func OpenDB() (*sql.DB, error) {
	return sql.Open("postgres", createDataSourceName())
}

// OpenGDB opens a postgres database specified by environment variable and wrapped by gorm.
// it shows log if it runs on local.
func OpenGDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", createDataSourceName())
	if err != nil {
		return nil, err
	}

	// callbacks for log
	db.Callback().Create().Register("log_sql", logSQLError)
	db.Callback().Update().Register("log_sql", logSQLError)
	db.Callback().Delete().Register("log_sql", logSQLError)
	db.Callback().Query().Register("log_sql", logSQLError)

	switch {
	case flag.Lookup("test.v") != nil:
		// no log in normal
		db.LogMode(false)
	case os.Getenv("RUNSERVER") == "LOCAL":
		// log on local
		db.LogMode(true)
	}
	return db, err
}

func createDataSourceName() string {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresName := os.Getenv("POSTGRES_NAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresName, postgresPassword)
}

// GetDB returns a pointer of database.
// if it hasn't set, it will return nil.
func GetDB() *gorm.DB {
	return db
}

// SetDB set a pointer of database.
func SetDB(d *gorm.DB) {
	db = d
}
