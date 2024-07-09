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

#  docker compose
dc-up:
	@docker compose up --build -d
dc-down:
	@docker compose down
