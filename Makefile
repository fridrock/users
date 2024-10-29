migrations:
	goose -dir db/migrations postgres "postgresql://users_service:root@127.0.0.1:5432/users?sslmode=disable" up
dropmigrations:
	goose -dir db/migrations postgres "postgresql://users_service:root@127.0.0.1:5432/users?sslmode=disable" down
run:
	go build -o users && ./users