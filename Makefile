.PHONY: list
list:
	@$(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | \
	awk -v RS= -F: '/^# Files/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") print $$1}' | \
	sort


# variable

# database migration
# database config
host=localhost
database=postgres
user=postgres
password=postgres
port=5432

# goose config
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgresql://$(user):$(password)@$(host):5432/$(database)?sslmode=disable
GOOSE_DIR=database/migrations

# create new migration with name
m-create:
	@goose -dir $(GOOSE_DIR) create $(name) sql
m-up:
	@goose -dir $(GOOSE_DIR) postgres $(GOOSE_DBSTRING) up
m-down:
	@goose -dir $(GOOSE_DIR) postgres $(GOOSE_DBSTRING) down
m-version:
	@goose -dir $(GOOSE_DIR) postgres $(GOOSE_DBSTRING) down-to $(version)
#  docker compose
dc-up:
	@docker compose up --build -d
dc-down:
	@docker compose down

# static check
lint:
	@staticcheck ./...

# swagger
s-fmt:
	@swag fmt
s-tidy: s-fmt
	@swag init -g ./cmd/app/api.go
