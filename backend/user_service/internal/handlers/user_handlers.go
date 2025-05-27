package user_handlers

import (
	"fmt"
	"net/http"
)

// func AddUser() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var
// 	}
// }

func Testfunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("user_handlers")
	}
}
