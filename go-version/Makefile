build:
	go build -o bin/api ./cmd/main.go

build-docker:
	@mkdir bin || true
	@docker run --rm --quiet -v $(PWD):/app -w /app golang:1.25.1 sh -c "go build -o bin/api ./cmd/main.go"

run:
	./bin/api

migrate-up:
	@docker run --rm --quiet  -v $(PWD):/app -w /app golang:1.25.1 sh -c "go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && migrate -path ./migrations -database 'sqlite://database.sqlite' up 1"

migrate-down:
	@docker run --rm --quiet  -v $(PWD):/app -w /app golang:1.25.1 sh -c "go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && migrate -path ./migrations -database 'sqlite://database.sqlite' down 1"

migrate-refresh: migrate-down migrate-up

test:
	@docker run --rm --quiet -v $(PWD):/app -w /app golang:1.25.1 sh -c "go test -v ./internal/api"

generate-mocks:
	mockgen -source=internal/api/store/users_store.go -destination=internal/api/store/mocks/mock_users_store.go -package=mocks
	mockgen -source=internal/api/store/reminders_store.go -destination=internal/api/store/mocks/mock_reminders_store.go -package=mocks

	mockgen -source=internal/api/repository/users_repository.go -destination=internal/api/repository/mocks/mock_users_repository.go -package=mocks
	mockgen -source=internal/api/repository/reminders_repository.go -destination=internal/api/repository/mocks/mock_reminders_repository.go -package=mocks