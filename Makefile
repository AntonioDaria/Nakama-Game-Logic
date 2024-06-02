start:
	@docker-compose up -d --build

stop:
	@docker-compose down

test:
	@go clean -testcache
	@go test -v ./...