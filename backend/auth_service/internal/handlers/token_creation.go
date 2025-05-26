package auth_handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

type UserTokenBody struct {
	UserId uint `json:"user_id"`
}

func CreateToken(redisDb *redis.Client, userId *uint) (string, error) {
	key, status := os.LookupEnv("TOKEN_SECRET_KEY")
	t := time.Now().Unix() + 60*60*24*31*6

	if !status {
		return "", fmt.Errorf("internal server error: secret key was not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     t,
	})

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("internal server error: %v", err.Error())
	}

	session_key := fmt.Sprintf("user:%v", strconv.Itoa(int(*userId)))

	_, err = redisDb.ZAdd(context.Background(), session_key, redis.Z{
		Score:  float64(t),
		Member: signed,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("internal server error: %v", err.Error())
	}

	return signed, nil
}

func CreateAdminToken(redisDb *redis.Client, userId *uint) (string, error) {
	key, status := os.LookupEnv("ADMIN_TOKEN_SECRET_KEY")
	t := time.Now().Unix() + 60*60*24*31*6

	if !status {
		return "", fmt.Errorf("internal server error: secret key was not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     t,
	})

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("internal server error: %v", err.Error())
	}

	session_key := fmt.Sprintf("admin:%v", strconv.Itoa(int(*userId)))

	_, err = redisDb.ZAdd(context.Background(), session_key, redis.Z{
		Score:  float64(t),
		Member: signed,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("internal server error: %v", err.Error())
	}

	return signed, nil
}
