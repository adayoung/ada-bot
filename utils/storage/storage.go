package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register DB driver for PostgreSQL
)

// DB is exported for use in sub-packages
var DB *sqlx.DB

// InitDB is called from main(), to set up DB and check for connectivity
func InitDB(connString string) error {
	var err error
	if DB, err = sqlx.Open("postgres", connString); err == nil {
		if err := DB.Ping(); err == nil {
			_onReady() // launch functions dependant on DBs readiness
		} else {
			return err // Error at DB.Ping() call
		}
	} else {
		return err // Error at sqlx.Open() call
	}
	return nil
}

var onReady []func() = []func(){}

// OnReady allows sub-packages to queue their own init() once we have a DB up
func OnReady(initdb func()) {
	onReady = append(onReady, initdb)
}

func _onReady() {
	for _, fn := range onReady {
		go fn() // launch functions dependant on DBs readiness
	}
}
