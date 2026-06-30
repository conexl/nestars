package telegram

import (
	"fmt"
	"net/http"
	"net/url"

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

	opts := []telego.BotOption{}

	if cfg.Telegram.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.Telegram.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse TELEGRAM_PROXY_URL: %w", err)
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}

		opts = append(opts, telego.WithHTTPClient(httpClient))
	}

	return telego.NewBot(cfg.Telegram.Token, opts...)
}
