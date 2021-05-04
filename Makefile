.PHONY: build 
build:
	@docker-compose build 

.PHONY: up
up:
	@docker-compose run --rm depends_resolver
	@docker-compose up -d app

.PHONY: down
down:
	@docker-compose down

.PHONY: exec/in
exec/in:
	@docker-compose exec $(CONTAINER) /bin/sh

.PHONY: in/app
in/app: CONTAINER=app
export CONTAINER
in/app:
	@$(MAKE) exec/in

.PHONY: in/db
in/db: CONTAINER=db
export CONTAINER
in/db:
	@$(MAKE) exec/in

.PHONY: reset/image
reset/image:
	@docker rmi voyagerwy130/sample_golang:1.0
	@docker rmi voyagerwy130/sample_db:1.0

.PHONY: reset/db
reset/db:
	@rm -rf db/data
	@mkdir db/data
