up-redis:
	docker-compose up -d redis

run:
	go run cmd/main.go

test:
	go test -v internal/infra/api/middleware/middleware_test.go