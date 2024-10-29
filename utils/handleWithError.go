package utils

import (
	"fmt"
	"net/http"
)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) (status int, err error)

func HandleErrorMiddleware(h HandlerWithError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := h(w, r)
		if err != nil {
			w.WriteHeader(status)
			w.Write([]byte(fmt.Sprint(err)))
		}
	})
}
