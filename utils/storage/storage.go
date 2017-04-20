package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register DB driver for PostgreSQL
)

// DB is exported for use in sub-packages
var DB *sqlx.DB

// Called from main(), sets up DB and checks for connectivity
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

// storage.OnReady allows sub-packages to queue their own init() sequentially
func OnReady(initdb func()) {
	onReady = append(onReady, initdb)
}

func _onReady() {
	for _, fn := range onReady {
		go fn() // launch functions dependant on DBs readiness
	}
}
