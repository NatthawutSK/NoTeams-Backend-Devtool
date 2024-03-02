DB_URL_PROD=
DB_URL_DEV=postgres://ri:123456@localhost:4444/noteam_devtool_db?sslmode=disable
PATH_MIGRATE ?= pkg/databases/migrations
TAG ?= v1

dev:
	air -c .air.dev.toml

prod:
	go run main.go .env.prod

init_db:
	docker run --name noteam_devtool_db -e POSTGRES_USER=ri -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:alpine

into_db:
	docker exec -it noteam_devtool_db bash -c 'psql -U ri'

create_db:
	docker exec -it noteam_devtool_db bash -c 'psql -U ri -c "CREATE DATABASE noteam_devtool_db;"'

drop_db:
	docker exec -it noteam_devtool_db bash -c 'psql -U ri -c "DROP DATABASE noteam_devtool_db;"'

db: init_db create_db

run_db:
	docker start noteam_devtool_db

migrate_up_prod:
	migrate -database '$(DB_URL_PROD)' -path $(PATH_MIGRATE) -verbose up

migrate_down_prod:
	migrate -database '$(DB_URL_PROD)' -path $(PATH_MIGRATE) -verbose down

migrate_up_dev:
	migrate -database '$(DB_URL_DEV)' -path $(PATH_MIGRATE) -verbose up

migrate_down_dev:
	migrate -database '$(DB_URL_DEV)' -path $(PATH_MIGRATE) -verbose down


build:
	docker build -t noteams-backend:$(TAG) .

docker_run:
	docker run -d -p 3000:3000 noteams-backend:$(TAG)

.PHONY: init_db into_db create_db drop_db db run_db migrate_up migrate_down dev prod into_db_prod build docker_run migrate_up_prod migrate_down_prod
