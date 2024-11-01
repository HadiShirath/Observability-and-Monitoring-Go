docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d

docker-reload:
	docker-compose stop microservice
	docker-compose build microservice --no-cache
	docker-compose start microservice

docker-down:
	docker-compose down 