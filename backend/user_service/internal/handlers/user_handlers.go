package user_handlers

import (
	"fmt"
	"net/http"
)

func Testfunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("user_handlers")
	}
}
