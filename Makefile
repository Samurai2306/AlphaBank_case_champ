.PHONY: api web seed test test-api up down

api:
	cd apps/api && go run ./cmd/api

seed:
	cd apps/api && go run ./cmd/seed

test-api:
	cd apps/api && go test ./...

web:
	cd apps/web && npm run dev

test:
	$(MAKE) test-api
	cd apps/web && npm test --if-present

up:
	docker compose up -d --build

down:
	docker compose down
