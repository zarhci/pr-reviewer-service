up-db:
	docker-compose up -d postgres

migrate-up:
	docker-compose run --rm migrate

goose-up:
	docker-compose run --rm goose

up-app:
	docker-compose up -d pr-reviewer

down:
	docker-compose down -v

up:
	docker-compose down -v
	docker-compose up -d postgres
	docker-compose run --rm migrate
	docker-compose up -d pr-reviewer

down:
	docker-compose down -v
