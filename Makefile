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
up:
	mkdir -p postgres_data
	docker-compose -f docker-compose-prod.yml up --build -d --remove-orphans
down:
	docker-compose -f docker-compose-prod.yml down
get.deps:
	rm -rf deps
	mkdir deps
	test -f deps/la/line-awesome.min.css || \
    	curl -o deps/la.zip https://maxst.icons8.com/vue-static/landings/line-awesome/line-awesome/1.3.0/line-awesome-1.3.0.zip && \
        unzip -d deps/la deps/la.zip;
	test -f deps/htmx.min.js || \
    	curl -o deps/htmx.min.js https://unpkg.com/htmx.org@2.0.1/dist/htmx.min.js;
	test -f deps/alpine.min.js || \
    	curl -o deps/alpine.min.js https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js;
	test -f deps/bulma.min.css || \
    	curl -o deps/bulma.min.css https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css;

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
redis:
	docker exec -it go-fiber-template-redis-1 redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT)
fmt:
	go fmt ./...
	TEMPL_EXPERIMENT=rawgo templ fmt .
dump.data:
	docker exec -it go-fiber-template-db-1 pg_dump $(DSN) --data-only -t contract_metadata -f /home/postgres/workdir/dump.data.sql
restore.data:
	docker exec -it go-fiber-template-db-1 psql $(DSN) -f /home/postgres/workdir/dump.data.sql
