.PHONY: test test-unit test-integration

test: test-unit test-integration

test-unit:
    go test ./tests/unit/... -v

test-integration:
    go test ./tests/integration/... -v