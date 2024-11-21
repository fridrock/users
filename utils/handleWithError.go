package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fridrock/users/api"
)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) (status int, err error)

func HandleErrorMiddleware(h HandlerWithError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := h(w, r)
		if err != nil {
			w.WriteHeader(status)
			slog.Error(err.Error())
			errMsg, err := json.MarshalIndent(api.ErrorResponse{
				Status:  status,
				Message: err.Error(),
			}, "", " ")
			if err != nil {
				slog.Error("Error encoding ErrorResponse")
			} else {
				w.Write([]byte(errMsg))
			}
		}
	})
}
