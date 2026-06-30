package telegram

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"go.uber.org/fx"
)

func Run(lc fx.Lifecycle, bot *telego.Bot, router *Router) {
	var cancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lpCtx, c := context.WithCancel(context.Background())
			cancel = c

			updates, err := bot.UpdatesViaLongPolling(lpCtx, nil)
			if err != nil {
				return fmt.Errorf("start long polling: %w", err)
			}

			go func() {
				defer cancel()
				for up := range updates {
					router.Dispatch(lpCtx, up)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})
}
