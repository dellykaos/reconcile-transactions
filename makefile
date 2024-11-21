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
