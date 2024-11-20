package database

import (
	"strings"
	"testing"

	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq" // Enable PostgreSQL driver if needed
	_ "modernc.org/sqlite"
)

func TestOpenWithUnsupportedDriver(t *testing.T) {
	db, err := Open("unsupported_driver", "", "", ":memory:", "", "")

	if err == nil {
		t.Fatal("err MUST NOT be nil")
	}

	if !strings.Contains(err.Error(), "driver unsupported_driver is not supported") {
		t.Fatal("err MUST contain 'unsupported_driver unsupported is not supported'")
	}

	if db != nil {
		t.Fatal("db MUST be nil")
	}
}

func TestOpen(t *testing.T) {
	db, err := Open(DRIVER_SQLITE, "", "", ":memory:", "", "")

	if err != nil {
		t.Fatal(err)
	}

	if db == nil {
		t.Fatal("db is nil")
	}
}
