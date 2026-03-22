APP=student-service-api
MIGRATIONS=./migrations
DEV_SEED_UP=./seeds/development/up.sql
DEV_SEED_DOWN=./seeds/development/down.sql
DB_URL?=$(DB_URL)
PSQL?=psql

.PHONY: run test tidy migrate-up migrate-down migrate-force seed-dev-up seed-dev-down bootstrap-dev

run:
	CGO_ENABLED=0 go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(VERSION)

seed-dev-up:
	$(PSQL) "$(DB_URL)" -v ON_ERROR_STOP=1 -f $(DEV_SEED_UP)

seed-dev-down:
	$(PSQL) "$(DB_URL)" -v ON_ERROR_STOP=1 -f $(DEV_SEED_DOWN)

bootstrap-dev:
	$(MAKE) migrate-up DB_URL="$(DB_URL)"
	$(MAKE) seed-dev-up DB_URL="$(DB_URL)" PSQL="$(PSQL)"
