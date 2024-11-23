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

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: run-reconcile-job
run-job:
	go run cmd/reconcile-job/main.go

.PHONY: build
build:
	go build -o output/api cmd/api/main.go

.PHONY: build-reconcile-job
build-reconcile-job:
	go build -o output/reconcile-job cmd/reconcile-job/main.go

.PHONY: compile.server
compile.server:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o deploy/api/api cmd/api/main.go
	docker buildx build -f dockerfile/api/Dockerfile --platform=linux/amd64 -t delly/amartha-recon-api:demo-amd64 .

.PHONY: compile.reconcile-job
compile.reconcile-job:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o deploy/reconcile-job/reconcile-job cmd/reconcile-job/main.go
	docker buildx build -f dockerfile/reconcile-job/Dockerfile --platform=linux/amd64 -t delly/amartha-recon-job:demo-amd64 .

.PHONY: test
test:
	go test -v ./...
