.PHONY: api web test test-api up down

api:
	cd apps/api && go run ./cmd/api

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
