run:
	go run ./cmd/main.go

build:
	go build -o bin/app ./cmd/main.go

migrate:
	go run ./cmd/migrate/main.go up

seed:
	go run ./cmd/migrate/main.go up

fe-build:
	cd frontend && npm run build

fe-dev:
	cd frontend && npm run dev

docker-up:
	docker compose up -d --build
