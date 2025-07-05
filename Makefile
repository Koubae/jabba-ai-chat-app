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
