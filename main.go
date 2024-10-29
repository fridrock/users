package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/auth_service/db/core"
	"github.com/fridrock/users/registration"
	"github.com/fridrock/users/utils"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	// http.ListenAndServe(":3000", r)
	startApp()
}

type App struct {
	server              *http.Server
	db                  *sqlx.DB
	registrationStore   registration.RegistrationStorage
	registrationHandler registration.RegistrationHandler
}

func startApp() {
	a := App{}
	a.setup()
}

func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.registrationStore = registration.NewRegistrationStorage(a.db)
	a.registrationHandler = registration.NewRegistrationHandler(a.registrationStore)
	a.setupServer()
}
func (a App) setupServer() {
	a.server = &http.Server{
		Addr:         ":9000",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      a.getRouter(),
	}
	slog.Info("Starting server on port 9000")
	a.server.ListenAndServe()
}
func (a App) getRouter() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/users/", utils.HandleErrorMiddleware(a.registrationHandler.HandleRegistration)).Methods("POST")
	return mainRouter
}

// func (a App) getUsersRouter(r *mux.Router) *mux.Router {
// 	usersRouter := r.PathPrefix("/users").Subrouter()
// 	usersRouter.Handle("/signup", handlers.HandleErrorMiddleware(a.userService.CreateUser)).Methods("POST")
// 	usersRouter.Handle("/send-confirmation", handlers.HandleErrorMiddleware(a.userService.SendConfirmation)).Methods("POST")
// 	usersRouter.Handle("/signin", handlers.HandleErrorMiddleware(a.userService.AuthUser)).Methods("POST")
// 	usersRouter.Handle("/confirm-email/{code}", handlers.HandleErrorMiddleware((a.userService.ConfirmEmail))).Methods("GET")
// 	return usersRouter
// }
