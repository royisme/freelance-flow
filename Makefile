GO ?= go
FRONTEND_DIR ?= frontend
PART ?= patch
VERSION ?=
NOTES ?=
OWNER ?= royzhu
REPO ?= freelance-flow
FILES ?=

.PHONY: lint test frontend-check build-darwin-amd64 build-darwin-arm64 build-darwin-all release-bump update-json migrate-up migrate-down migrate-steps migrate-version migrate-force ent-generate

lint:
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint run

test:
	$(GO) test ./internal/...

frontend-check:
	cd $(FRONTEND_DIR) && bun install && bun run check

build-darwin-amd64:
	wails build -platform darwin/amd64 -ldflags "-X freelance-flow/internal/update.CurrentVersion=$$(git describe --tags --always)"

build-darwin-arm64:
	wails build -platform darwin/arm64 -ldflags "-X freelance-flow/internal/update.CurrentVersion=$$(git describe --tags --always)"

build-darwin-arm64-no-finance:
	FINANCE_MODULE_DISABLED=1 wails build -platform darwin/arm64 -tags FINANCE_MODULE_DISABLED -ldflags "-X freelance-flow/internal/update.CurrentVersion=$$(git describe --tags --always)"

build-darwin-all: build-darwin-amd64 build-darwin-arm64

release-bump:
	$(GO) run ./cmd/release/main.go -part $(PART)

update-json:
	@if [ -z "$(VERSION)" ] || [ -z "$(FILES)" ]; then \
		echo "Usage: make update-json VERSION=1.2.3 NOTES='...' FILES='path1 path2' [OWNER=royzhu REPO=freelance-flow OUTPUT=update.json]"; \
		exit 1; \
	fi
	$(GO) run ./cmd/gen-update-info/main.go -version $(VERSION) -notes "$(NOTES)" -owner $(OWNER) -repo $(REPO) -output update.json $(FILES)

migrate-up:
	$(GO) run ./cmd/migrate/main.go -action up

migrate-down:
	$(GO) run ./cmd/migrate/main.go -action down

migrate-steps:
	$(GO) run ./cmd/migrate/main.go -action steps -steps $(STEPS)

migrate-version:
	$(GO) run ./cmd/migrate/main.go -action version

migrate-force:
	$(GO) run ./cmd/migrate/main.go -action force -version $(VERSION)

ent-generate:
	$(GO) run entgo.io/ent/cmd/ent generate ./ent/schema

