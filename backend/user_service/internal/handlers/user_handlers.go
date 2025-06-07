package user_handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func CheckUserExistence(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("asdf")
	}
}
