new-migration-file:
	migrate create -ext sql -dir internal/db/migrations $(name)

sqlc:
	sqlc generate

grey-app-db-debug:
	docker run --name grey-app-db-debug -p 5444:5432 -e POSTGRES_PASSWORD=db_pass -e POSTGRES_USER=db_user -e POSTGRES_DB=grey-app-db -d postgres

stop-grey-app-db-debug:
	docker stop grey-app-db-debug

remove-grey-app-db-debug: stop-grey-app-db-debug
	docker rm grey-app-db-debug

build-app-binary:
	sqlc generate
	GOOS=linux GOARCH=amd64 go build -ldflags=-s -o main cmd/wallet-app/main.go

start-wallet:
	docker compose --env-file dev.env up app

start-all-services: build-app-binary
	docker compose --env-file dev.env up --build -d

start-all-services-and-seed-dev: start-all-services
	MY_ENV=development go run cmd/seeder/main.go

stop-all-services:
	docker compose --env-file dev.env stop

kill-all-services-with-all-data:
	docker compose --env-file dev.env down -v
	docker rmi grey-app

reload: kill-all-services-with-all-data start-all-services

app-logs:
	docker logs grey-wallet-backend-app