DB_URL_PROD=postgres://noteams:njDcP7FR2YEy3yFk@noteams-prod.c1qaumiuo6j0.us-east-1.rds.amazonaws.com:5432/noteams_db?sslmode=disable
DB_URL_DEV=postgres://ri:123456@localhost:4444/noteams_db_dev?sslmode=disable
PATH_MIGRATE ?= pkg/databases/migrations
TAG ?= v1

dev:
	air -c .air.dev.toml

prod:
	go run main.go .env.prod

init_db:
	docker run --name noteams_db_dev -e POSTGRES_USER=ri -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:alpine

into_db:
	docker exec -it noteams_db_dev bash -c 'psql -U ri'

create_db:
	docker exec -it noteams_db_dev bash -c 'psql -U ri -c "CREATE DATABASE noteams_db_dev;"'

drop_db:
	docker exec -it noteams_db_dev bash -c 'psql -U ri -c "DROP DATABASE noteams_db_dev;"'

db: init_db create_db

run_db:
	docker start noteams_db_dev

migrate_up_prod:
	migrate -database '$(DB_URL_PROD)' -path $(PATH_MIGRATE) -verbose up

migrate_down_prod:
	migrate -database '$(DB_URL_PROD)' -path $(PATH_MIGRATE) -verbose down

migrate_up_dev:
	migrate -database '$(DB_URL_DEV)' -path $(PATH_MIGRATE) -verbose up

migrate_down_dev:
	migrate -database '$(DB_URL_DEV)' -path $(PATH_MIGRATE) -verbose down

into_db_prod:
	psql --host=noteams-prod.c1qaumiuo6j0.us-east-1.rds.amazonaws.com --port=5432 --username=noteams --password --dbname=noteams_db


clone_git: 
	git clone https://NatthawutSK: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@github.com/NatthawutSK/NoTeams-Backend.git

build:
	docker build -t noteams-backend:$(TAG) .

docker_run:
	docker run -d -p 3000:3000 noteams-backend:$(TAG)

.PHONY: init_db into_db create_db drop_db db run_db migrate_up migrate_down dev prod into_db_prod clone_git build docker_run migrate_up_prod migrate_down_prod
# git clone https://NatthawutSK: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@github.com/NatthawutSK/NoTeams-Backend.git
# https://github.com/NatthawutSK/NoTeams-Backend.git
