package telegram

import (
	"fmt"

	"github.com/mymmrac/telego"
	"go.uber.org/fx"

	"nestars-bot/internal/config"
	"nestars-bot/internal/infrastructure/telegram/handlers"
)

var Module = fx.Module("telegram",
	fx.Provide(
		NewBot,
		NewFSM,
		NewRender,
		// Адаптируем конкретные типы под интерфейсы хендлеров.
		func(r *Render) handlers.Render { return r },
		func(f *FSM) handlers.FSM { return f },
		handlers.NewMenu,
		handlers.NewProfile,
		handlers.NewStars,
		handlers.NewPremium,
		handlers.NewPayment,
		NewRouter,
	),
	fx.Invoke(Run),
)

func NewBot(cfg *config.Config) (*telego.Bot, error) {
	if cfg.Telegram.Token == "" {
		return nil, fmt.Errorf("BOT_TOKEN is not set")
	}
	return telego.NewBot(cfg.Telegram.Token)
}
