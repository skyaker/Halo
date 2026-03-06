package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	handlers "auth_service/internal/handlers"
	dbconn "auth_service/internal/repository"
	service "auth_service/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func main() {
	// postgres
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	// redis
	redisDb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	_, err := redisDb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth service").
			Msg("Redis connection failed")
	}
	log.Info().Msg("Redis connection successful")

	go listenForExpiredTokens(redisDb)

	// kafka
	kafkaUrl := fmt.Sprintf("%v:%v", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	writer := getKafkaWriter(kafkaUrl)

	log.Info().Msg("Kafka writer created")

	defer writer.Close()

	// token secret key
	key, status := os.LookupEnv("TOKEN_SECRET_KEY")
	if !status {
		log.Fatal().Msg("TOKEN_SECRET_KEY environment variable is not set")
	}

	// router
	r := chi.NewRouter()

	authService := service.NewAuthService(db, redisDb, writer, key)
	authHandler := handlers.NewAuthHandler(authService)

	r.Post("/api/auth/register", authHandler.HandleRegister())
	r.Delete("/api/auth/delete_user", authHandler.HandleDelete())
	r.Get("/api/auth/check_token", authHandler.CheckToken())
	r.Post("/api/auth/login", authHandler.Login())

	log.Info().Msg("Auth server is running")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth service").
			Msg("Server start failed")
	}
}

func getKafkaWriter(kafkaURL string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(kafkaURL),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}
}

func listenForExpiredTokens(redisDb *redis.Client) {
	ctx := context.Background()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()
		users, err := redisDb.Keys(ctx, "user:*").Result()
		if err != nil {
			log.Error().Err(err).Msg("Error fetching keys")
			continue
		}

		for _, userKey := range users {
			removed, err := redisDb.ZRemRangeByScore(ctx, userKey, "0", fmt.Sprintf("%d", now)).
				Result()
			if err != nil {
				log.Error().Err(err).Msg("Error moving expired tokens")
				continue
			}

			if removed > 0 {
				log.Info().Msgf("Removed %d expired tokens from %s\n", removed, userKey)
			}
		}
	}
}
