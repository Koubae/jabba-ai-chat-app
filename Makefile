.PHONY: run build stop tests

# ============================
# 	Run
# ============================

# //////////////////////
# 	chat-identity
# //////////////////////

run-chat-identity:
	@air -c .air.chat-identity.toml

run-chat-identity-local:
	go run cmd/chat-identity/main.go

# //////////////////////
# 	chat-orchestrator
# //////////////////////

run-chat-orchestrator:
	@air -c .air.chat-orchestrator.toml

run-chat-orchestrator-local:
	go run cmd/chat-orchestrator/main.go

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
