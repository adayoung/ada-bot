package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(conn_string string) error {
	var err error
	if DB, err = sqlx.Open("postgres", conn_string); err == nil {
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

func OnReady(initdb func()) {
	onReady = append(onReady, initdb)
}

func _onReady() {
	for _, fn := range onReady {
		go fn() // launch functions dependant on DBs readiness
	}
}
