module github.com/fridrock/users

go 1.22.1

require (
	github.com/fridrock/auth_service v0.0.0-20240503143716-21dc4136010f
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.28.0
)

require github.com/gorilla/mux v1.8.1

require github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
