.PHONY: dev test test-backend test-frontend integration-test build lint docker-up docker-down migrate seed clean

dev:
	docker compose up --build
test: test-backend test-frontend
test-backend:
	cd backend && go test ./...
test-frontend:
	cd frontend && npm test
integration-test:
	cd backend && go test -count=1 ./tests/integration
build:
	cd backend && go build ./cmd/api
	cd frontend && npm run build
lint:
	cd backend && test -z "$$(gofmt -l .)" && go vet ./...
	cd frontend && npm run lint
docker-up:
	docker compose up --build -d
docker-down:
	docker compose down
migrate:
	docker compose run --rm migrate
seed:
	docker compose exec -T mysql sh -c 'mysql -u"$$MYSQL_USER" -p"$$MYSQL_PASSWORD" "$$MYSQL_DATABASE"' < scripts/seed.sql
clean:
	rm -rf frontend/dist frontend/coverage frontend/playwright-report

