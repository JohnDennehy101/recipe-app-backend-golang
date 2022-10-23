package driver

import (
	"database/sql"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var databaseConnection = &DB{}

const maxOpenDatabaseConnections = 10
const maxIdleDatabaseConnections = 5
const maxDatabaseLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDatabaseConnections)
	d.SetMaxIdleConns(maxIdleDatabaseConnections)
	d.SetConnMaxLifetime(maxDatabaseLifetime)

	databaseConnection.SQL = d
	err = testDatabaseConnection(d)

	if err != nil {
		return nil, err
	}

	return databaseConnection, nil
}

func testDatabaseConnection(d *sql.DB) error {
	err := d.Ping()

	if err != nil {
		return err
	}

	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
