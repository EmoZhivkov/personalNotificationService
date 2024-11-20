start-all:
	docker-compose up --build

stop-all:
	docker-compose down

test:
	CGO_ENABLED=1 go test ./...

start-postgres:
	docker-compose up postgres

stop-postgres:
	docker-compose down postgres

generate-mocks:
	mockery --all

generate-swagger-docs:
	swag init

.PHONY: start-all stop-all test docker-run docker-build generate-mocks generate-swagger-docs
