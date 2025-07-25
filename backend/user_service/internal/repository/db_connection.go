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
	host, envSt := os.LookupEnv("USER_POSTGRES")
	if !envSt {
		log.Fatal().Msg("User host name not found")
		return nil
	}

	port, envSt := os.LookupEnv("USER_PG_PORT")
	if !envSt {
		log.Fatal().Msg("User postgres port not found")
		return nil
	}

	portInt, _ := strconv.Atoi(port)
	portUint := uint(portInt)

	user, envSt := os.LookupEnv("POSTGRES_USER")
	if !envSt {
		log.Fatal().Msg("User postgres user not found")
		return nil
	}

	password, envSt := os.LookupEnv("POSTGRES_PASSWORD")
	if !envSt {
		log.Fatal().Msg("User postgres password not found")
		return nil
	}

	dbname, envSt := os.LookupEnv("USER_DB")
	if !envSt {
		log.Fatal().Msg("User db name not found")
		return nil
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
			Msg("pg open")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("pg ping")
	}

	log.Info().
		Str("service", "user_service").
		Msg("Postgres connection successfull")

	return db
}
