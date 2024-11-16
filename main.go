package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/auth_service/db/core"
	"github.com/fridrock/users/friend"
	"github.com/fridrock/users/usr"
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
	userStorage         usr.UserStorage
	userHandler         usr.UserHandler
	friendStorage       friend.FriendStorage
	friendHandler       friend.FriendHandler
	tokenRefreshHandler usr.TokenRefreshHandler
	authManager         utils.AuthManager
}

func startApp() {
	a := App{}
	a.setup()
}

func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.userStorage = usr.NewUserStorage(a.db)
	a.userHandler = usr.NewUserHandler(a.userStorage)
	a.friendStorage = friend.NewFriendStorage(a.db)
	a.friendHandler = friend.NewFriendHandler(a.friendStorage)
	a.tokenRefreshHandler = usr.NewTokenRefreshHandler()
	a.authManager = utils.NewAuthManager()
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
	mainRouter.Handle("/users/reg", utils.HandleErrorMiddleware(a.userHandler.HandleRegistration)).Methods("POST")
	mainRouter.Handle("/users/auth", utils.HandleErrorMiddleware(a.userHandler.HandleAuth)).Methods("POST")
	mainRouter.Handle("/users/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.userHandler.FindUser))).Methods("GET")
	mainRouter.Handle("/token/refresh", utils.HandleErrorMiddleware((a.tokenRefreshHandler.HandleRefreshToken))).Methods("POST")
	mainRouter.Handle("/friends/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.friendHandler.AddFriend))).Methods("POST")
	mainRouter.Handle("/friends/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.friendHandler.DeleteFriend))).Methods("DELETE")
	mainRouter.Handle("/friends/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.friendHandler.GetFriends))).Methods("GET")
	return mainRouter
}
