APP_NAME=product-command
GO_CMD=go
ATLAS=atlas

run:
	$(GO_CMD) run cmd/app/main.go

migrate.diff:
	$(ATLAS) migrate diff --env local

migrate:
	$(ATLAS) migrate apply --env local

migrate.status:
	$(ATLAS) migrate status --env local

migrate.lint:
	$(ATLAS) migrate lint --env local

migrate.hash:
	$(ATLAS) migrate hash --env local

fmt:
	$(GO_CMD) fmt ./...