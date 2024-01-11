.PHONY: init

init:
	@echo "---init---"
	brew install pre-commit
	brew install golangci-lint
	# @echo "== pre-commit setup =="
	pre-commit install

tidy:
	go mod tidy

mocks:
	@echo "Generate mock repository"
	cd app/repository && mockery --all --case underscore --output repository_mock --outpkg repository_mock

run-linter:
	@echo Starting linters
	golangci-lint run ./...

test:
	@echo Starting run unittest
	go test ./...

precommit.rehooks:
	pre-commit autoupdate
	pre-commit install --install-hooks
	pre-commit install --hook-type prepare-commit-msg
	pre-commit install --hook-type pre-commit
	pre-commit install --hook-type commit-msg
	pre-commit install --hook-type pre-push
	pre-commit install --hook-type post-commit

precommit.clean:
	pre-commit clean

precommit.gc:
	pre-commit gc

run-server:
	go run main.go

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

dev:
	docker-compose -f docker-compose.dev.yml up

scan-vulnerability:
	osv-scanner --lockfile './go.mod'

scan-security:
	gosec -exclude=G302,G304,G505,G107,G404 ./...
