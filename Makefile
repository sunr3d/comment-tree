up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

clean:
	docker compose down -v

logs:
	docker compose logs -f app

test:
	go test -v ./...

migrate-up:
	docker compose exec -T db psql -U comment_tree_user -d comment_tree -f /migrations/001_init_up.sql

migrate-down:
	docker compose exec -T db psql -U comment_tree_user -d comment_tree -f /migrations/001_init_down.sql

fmt:
	go fmt ./...

build:
	go build -o comment-tree cmd/main.go