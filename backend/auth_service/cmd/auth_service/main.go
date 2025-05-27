package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handlers "auth_service/internal/handlers"
	dbconn "auth_service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	var db *sql.DB = dbconn.GetDbConnection()
	defer db.Close()

	redisDb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	_, err := redisDb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	go listenForExpiredTokens(redisDb)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		MaxAge:           300,
	}))

	r.Post("/api/auth/register", handlers.RegisterUser(db, redisDb))
	r.Post("/api/auth/check_token", handlers.CheckToken(db, redisDb))
	r.Post("/api/auth/login", handlers.Login(db, redisDb))

	log.Println("Auth server is running")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server failed to start:", err)
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
			log.Println("Error fetching keys:", err)
			continue
		}

		for _, userKey := range users {
			removed, err := redisDb.ZRemRangeByScore(ctx, userKey, "0", fmt.Sprintf("%d", now)).
				Result()
			if err != nil {
				log.Println("Error removing expired tokens:", err)
				continue
			}

			if removed > 0 {
				fmt.Printf("Removed %d expired tokens from %s\n", removed, userKey)
			}
		}
	}
}
