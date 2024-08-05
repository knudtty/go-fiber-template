include $(shell echo .env)
export DB_USER
export DB_PASSWORD
export DB_NAME

TEMPL_VERSION     := v0.2.747
export TEMPL_VERSION
BULMA_VERSION     := 1.0.2
export BULMA_VERSION
ifeq ($(SERVER_NAME),localhost)
	NGINX_CONF := ./nginx/developent.conf.template
else
	NGINX_CONF := ./nginx/production.conf.template
endif

watch:
	docker-compose -f docker-compose-dev.yml up --build --remove-orphans
watch-down:
	docker-compose -f docker-compose-dev.yml down
up:
	docker-compose up --build -d --remove-orphans
down:
	docker-compose down
update:
	go get -u github.com/a-h/templ
	go get -u github.com/gofiber/fiber/v3
	go mod tidy
