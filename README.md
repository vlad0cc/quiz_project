# quiz_project
small quiz for university as pet project 

Backend - Golang 1.22.0

Frontend - Typescript 6.0 ; React; HTML; CSS

DB - PostgreSQL 16


Commands to run the project in order:

docker-up:
	docker compose up -d --build

migrate:
	go run ./cmd/migrate/main.go up

fe-build:
	cd frontend && npm run build

fe-dev:
	cd frontend && npm run dev

seed:
	go run ./cmd/migrate/main.go up

build:
	go build -o bin/app ./cmd/main.go

run:
	go run ./cmd/main.go
