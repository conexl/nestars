package main

import (
	"go.uber.org/fx"

	"nestars-bot/internal/config"
	"nestars-bot/internal/infrastructure/api"
	"nestars-bot/internal/infrastructure/miniapp"
	"nestars-bot/internal/infrastructure/persistence"
	"nestars-bot/internal/infrastructure/telegram"
	"nestars-bot/internal/usecase"
)

func main() {
	fx.New(
		fx.Provide(config.NewConfig),
		persistence.Module,
		api.Module,
		telegram.Module,
		miniapp.Module,
		usecase.Module,
	).Run()
}
