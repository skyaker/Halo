package connection

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func GetDbConnection() *sql.DB {
	host, _ := os.LookupEnv("AUTH_POSTGRES")

	port, _ := os.LookupEnv("AUTH_PG_PORT")
	portInt, _ := strconv.Atoi(port)
	portUint := uint(portInt)

	user, _ := os.LookupEnv("POSTGRES_USER")
	password, _ := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, _ := os.LookupEnv("AUTH_DB")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host,
		portUint,
		user,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth_service").
			Str("process", "pg open")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth_service").
			Str("process", "pg ping")
	}

	return db
}
