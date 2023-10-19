up:
	docker compose -f docker-compose.yml  up -d --build
stop:
	docker compose stop logging_bot
start:
	docker compose start logging_bot
