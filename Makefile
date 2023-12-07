include .env

run: .check-env-vars
	docker-compose up

.check-env-vars:
	@test $${APOD_API?Please set environment variable APOD_API}