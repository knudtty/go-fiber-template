include $(shell echo .env)
TEMPL_VERSION     := v0.2.747
BULMA_VERSION     := 1.0.2
export DB_USER
export DB_PASSWORD
export DB_NAME
export TEMPL_VERSION
export BULMA_VERSION

watch:
	docker-compose -f docker-compose-dev.yml up --build --remove-orphans
watch.down:
	docker-compose -f docker-compose-dev.yml down
update:
	go get -u
	go mod tidy

DSN:=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
migrate.up:
	docker run --rm -it -v ./platform/migrations:/migrations migrate/migrate -path=/migrations/ -database $(DSN) up
migrate.down:
	docker run --rm -it -v ./platform/migrations:/migrations migrate/migrate -path=/migrations/ -database $(DSN) down
psql:
	docker exec -it go-fiber-template-db-1 psql $(DSN)
fmt:
	go fmt ./...
	TEMPL_EXPERIMENT=rawgo templ fmt .
