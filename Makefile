up-backend:
	docker compose up backend backend-ro-1 backend-ro-2 backend-mirror postgres-master postgres-slave-ro -d

up-db:
	docker compose up postgres-master postgres-slave-ro -d

up-web-server:
	docker compose up pgadmin swagger nginx -d

up-load-testing-gatling:
	docker compose up load-testing-gatling -d

build-backend:
	cd backend && go mod tidy
	cd backend && go build src\main.go

test-unit:
	rm -rf ./backend/src/tests/unit/allure-results

	cd backend && go test ./src/tests/unit/* --race --parallel 8

	mkdir ./backend/src/tests/unit/allure-results
	mv ./backend/src/tests/unit/pkg/allure-results/* ./backend/src/tests/unit/allure-results
	mv ./backend/src/tests/unit/service/allure-results/* ./backend/src/tests/unit/allure-results
	mv ./backend/src/tests/unit/repository/allure-results/* ./backend/src/tests/unit/allure-results

	rm -rf ./backend/src/tests/unit/pkg/allure-results
	rm -rf ./backend/src/tests/unit/repository/allure-results
	rm -rf ./backend/src/tests/unit/service/allure-results

up-allure-report:
	docker compose up allure allure-ui -d

.PHONY: up-backend up-db up-web-server up-load-testing-gatling build-backend test-unit up-allure-report