package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handlers "auth_service/internal/handlers"
	middleware "auth_service/internal/middleware"
	dbconn "auth_service/internal/repository"
	service "auth_service/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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

	// kafka
	kafkaUrl := fmt.Sprintf("%v:%v", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	writer := service.GetKafkaWriter(kafkaUrl)

	log.Info().Msg("Kafka writer created")

	defer writer.Close()

	// token secret key
	key, status := os.LookupEnv("TOKEN_SECRET_KEY")
	if !status {
		log.Fatal().Msg("TOKEN_SECRET_KEY environment variable is not set")
	}

	// session prefix
	sessionPrefix, status := os.LookupEnv("REDIS_USER_PREFIX")
	if !status {
		log.Fatal().Msg("SESSION_PREFIX environment variable is not set")
	}

	// router
	r := chi.NewRouter()

	authService := service.NewAuthService(db, redisDb, writer, key, sessionPrefix)
	authHandler := handlers.NewAuthHandler(authService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go authService.StartTokenCleanup(ctx)

	r.Post("/api/auth/register", authHandler.HandleRegister())
	r.Post("/api/auth/login", authHandler.Login())

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))
		r.Delete("/api/auth/delete_user", authHandler.HandleDelete())
		r.Get("/api/auth/me", authHandler.HandleMe())
	})

	log.Info().Msg("Auth server is running")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "auth service").
			Msg("Server start failed")
	}
}
