package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the db connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}


const maxOpenDbConns = 10
const maxIdleDbConns = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL connects database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConns)
	d.SetMaxIdleConns(maxIdleDbConns)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	if err := testDB(d); err != nil {
		return nil, err
	}

	return dbConn, nil
}

func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}