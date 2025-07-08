package connection

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/rs/zerolog/log"
)

func GetDbConnection() *sql.DB {
	host, envSt := os.LookupEnv("AUTH_POSTGRES")
	if !envSt {
		log.Fatal().Msg("Auth host name not found")
	}

	port, envSt := os.LookupEnv("AUTH_PG_PORT")
	if !envSt {
		log.Fatal().Msg("Auth postgres porn not found")
	}

	portInt, _ := strconv.Atoi(port)
	portUint := uint(portInt)

	user, envSt := os.LookupEnv("POSTGRES_USER")
	if !envSt {
		log.Fatal().Msg("Auth postgres user not found")
		return nil
	}

	password, envSt := os.LookupEnv("POSTGRES_PASSWORD")
	if !envSt {
		log.Fatal().Msg("Auth postgres password not found")
		return nil
	}

	dbname, envSt := os.LookupEnv("AUTH_DB")
	if !envSt {
		log.Fatal().Msg("Auth db name not found")
	}

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

	log.Info().
		Str("service", "auth_service").
		Msg("Postgres connection successfull")

	return db
}
