run:
	docker-compose -f docker/dev/docker-compose.yaml up --build

#default value for TAG
TAG ?= prod-img:latest
build:
	docker build -t $(TAG) -f docker/dev/Dockerfile .

run-prod:
	docker-compose -f docker/prod/docker-compose.yaml up --build
