package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/timelyrain/star-account/api"
	"github.com/timelyrain/star-account/db"
)

const (
	dbDriver          = "postgres"
	dbSource          = "postgresql://root:secret@localhost:15432/star_account?sslmode=disable"
	serverAddress     = "0.0.0.0:20714"
	tokenSymmetricKey = "11451419198101145141919810aaaaaa"
	tokenDuration     = time.Hour
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	db := db.NewDB(conn)

	var server *api.Server
	server, err = api.NewServer(db, tokenSymmetricKey, tokenDuration)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
