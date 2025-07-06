.PHONY: run build stop tests

# ============================
# 	Run
# ============================

# //////////////////////
# 	chat-identity
# //////////////////////

chat-identity-run:
	@air -c .air.chat-identity.toml

chat-identity-run-local:
	go run cmd/chat-identity/main.go

# //////////////////////
# 	chat-orchestrator
# //////////////////////

chat-orchestrator-run:
	@air -c .air.chat-orchestrator.toml

chat-orchestrator-run-local:
	go run cmd/chat-orchestrator/main.go

# //////////////////////
# 	chat-session
# //////////////////////

chat-session-run:
	@air -c .air.chat-session.toml

chat-session-run-local:
	go run cmd/chat-session/main.go


# ============================
#       Docker
# ============================
up:
	@docker compose up

down:
	@docker compose down

down-v:
	@docker compose down -v

build:
	@echo 'Building images ...🛠️'
	@docker compose build


# //////////////////////
# 	chat-identity
# //////////////////////
up-chat-identity:
	@docker compose up chat-identity

down-chat-identity:
	@docker compose down chat-identity

down-chat-identity-v:
	@docker compose down chat-identity -v


# //////////////////////
# 	chat-orchestrator
# //////////////////////
up-chat-orchestrator:
	@docker compose up chat-orchestrator

down-chat-orchestrator:
	@docker compose down chat-orchestrator

down-chat-orchestrator-v:
	@docker compose down chat-orchestrator -v


# //////////////////////
# 	chat-session
# //////////////////////
up-chat-session:
	@docker compose up chat-session

down-chat-session:
	@docker compose down chat-session

down-chat-session-v:
	@docker compose down chat-session -v

# ============================
# 	Tests
# ============================

### Using locally installed Go
test-unit-local:
	go test -v ./... -coverprofile=coverage.out && go tool cover -o coverage.html -html=coverage.out

test-specific:
ifndef TEST
	@echo "Please provide a test pattern using TEST=<pattern>"
	@echo "Example: make test-specific TEST=TestGetEnv/string_tests"
	@echo "make test-specific TEST=TestGetEnv"
	@echo "make test-specific TEST=TestGetEnv/string"
	@echo "make test-specific TEST=TestGetEnv/int"
	@echo "make test-specific TEST=TestGetEnv/int"
	@echo "make test-specific TEST=TestGetEnv/int_tests"
	@echo "\nAvailable test patterns:"
	@go test ./... -v -list=. | grep "^Test"
else
	go test ./... -v -run $(TEST)
endif


# --------------------------
# Init
# --------------------------
init: .install-deps

.install-deps:
	go mod tidy



# ============================
# 	Management
# ============================
script-prepare-mongo-database:
	go run scripts/database/prepare_mongo_database.go
