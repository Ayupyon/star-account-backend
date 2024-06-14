package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/timelyrain/star-account/db/sqlc"
)

const (
	testDbDriver = "postgres"
	testDbSource = "postgresql://root:secret@localhost:15432/star_account?sslmode=disable"
)

var testQueries *sqlc.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(testDbDriver, testDbSource)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	testQueries = sqlc.New(conn)

	os.Exit(m.Run())
}
