.PHONY: up down api miniapp build-api build-miniapp build-bot run-bot bot logs

# Run API + miniapp on the current server (no Telegram bot)
up:
	docker compose up -d --build api miniapp

down:
	docker compose down

api:
	docker compose up -d --build api

miniapp:
	docker compose up -d --build miniapp

build-api:
	docker compose build api

build-miniapp:
	docker compose build miniapp

# Build Docker image for the Telegram bot (API/miniapp excluded)
build-bot:
	docker build --target bot -t nestars-bot -f ./bot/Dockerfile .

# Run the Telegram bot in Docker (API/miniapp excluded)
run-bot:
	docker run --rm --env-file ./bot.env nestars-bot

# Run the Telegram bot locally with Go (API/miniapp excluded)
bot:
	go run ./bot/cmd/bot/main.go

logs:
	docker compose logs -f
