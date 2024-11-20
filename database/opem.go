package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

func Options() *openOptions {
	return &openOptions{
		properties: make(map[string]interface{}),
	}
}

type openOptions struct {
	properties map[string]interface{}
}

// Open opens the database
//
// Note:
//   - drivers are not included to this package to prevent size bloat
//   - you must add only the required database driver
//
// Drivers:
// - sqlite add the following includes:
// ```
// _ "modernc.org/sqlite"
// ```
// - mysql add the following includes:
// ```
// _ "github.com/go-sql-driver/mysql"
// ```
// - postgres add the following includes:
// ```
// _ "github.com/lib/pq"
// ```
//
// Business logic:
//   - opens the database based on the driver name
//   - each driver has its own set of parameters
//
// Parameters:
// - driverName: the driver name
// - dbHost: the database host
// - dbPort: the database port
// - dbName: the database name
// - dbUser: the database user
// - dbPass: the database password
//
// Returns:
// - *sql.DB: the database connection
// - error: the error if any
func Open(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	supportedDrivers := []string{DRIVER_SQLITE, DRIVER_MYSQL, DRIVER_POSTGRES}

	if !strings.EqualFold(driverName, DRIVER_SQLITE) &&
		!strings.EqualFold(driverName, DRIVER_MYSQL) &&
		!strings.EqualFold(driverName, DRIVER_POSTGRES) {
		return nil, errors.New(`driver ` + driverName + ` is not supported. Supported drivers: ` + strings.Join(supportedDrivers, ", "))
	}

	if strings.EqualFold(driverName, DRIVER_SQLITE) {
		dsn := dbName
		db, err = sql.Open(DRIVER_SQLITE, dsn)
	}

	if strings.EqualFold(driverName, DRIVER_MYSQL) {
		dsn := dbUser + `:` + dbPass + `@tcp(` + dbHost + `:` + dbPort + `)/` + dbName + `?charset=utf8mb4&parseTime=True&loc=UTC`
		db, err = sql.Open(DRIVER_MYSQL, dsn)
		// Maximum Idle Connections
		db.SetMaxIdleConns(5)
		// Maximum Open Connections
		db.SetMaxOpenConns(5)
		// Idle Connection Timeout
		db.SetConnMaxIdleTime(5 * time.Second)
		// Connection Lifetime
		db.SetConnMaxLifetime(30 * time.Second)
	}

	if strings.EqualFold(driverName, DRIVER_POSTGRES) {
		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
		db, err = sql.Open(DRIVER_POSTGRES, dsn)
	}

	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, errors.New("database for driver " + driverName + " could not be intialized")
	}

	err = db.Ping()

	if err != nil {
		return nil, errors.Join(errors.New("database for driver "+driverName+" could not be pinged"), err)
	}

	return db, nil
}
