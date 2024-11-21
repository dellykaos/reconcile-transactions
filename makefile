include .env

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir db/migrations $(name)

.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down 1

.PHONY: migrate-force
migrate-force:
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose force $(version)

.PHONY: setup-db
setup-db:
	@echo "Creating database '$(DB_NAME)'..."
	@PGPASSWORD=$(DB_PASS) psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);" || echo "Database '$(DB_NAME)' already exists."
	@make migrate-up

.PHONY: drop-db
drop-db:
	@echo "Dropping database '$(DB_NAME)'..."
	@PGPASSWORD=$(DB_PASS) psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE $(DB_NAME);" || echo "Database '$(DB_NAME)' does not exist."
