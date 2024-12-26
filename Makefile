run-docker-compose:
	docker compose up -d

build-backend:
	cd backend && go mod tidy
	cd backend && go build src\main.go

.PHONY:  build-app