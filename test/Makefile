.PHONY: test contract fuzz lint

# Run all tests
test:
    go test ./... -v

# Run contract tests (dynamic endpoint testing)
contract:
    cd test/contract && go test -run Test_OpenAPI_Contract -v

# Run authentication flow tests
contract-auth:
    cd test/contract && go test -run Test_OpenAPI_Contract_AuthFlow -v

# Run all contract tests
contract-all:
    cd test/contract && go test -v

# Validate OpenAPI spec
validate-spec:
    cd test/contract && go test -run ValidateSpecFile -v

# Run tests from project root (alternative)
contract-root:
    go test ./test/contract -run Test_OpenAPI_Contract -v

# Run fuzz testing
fuzz:
    docker compose -f docker-compose.schemathesis.yaml up --abort-on-container-exit --exit-code-from schemathesis

# Lint code
lint:
    golangci-lint run