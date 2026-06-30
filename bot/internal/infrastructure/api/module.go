package api

import (
	"nestars-bot/internal/config"
	"nestars-bot/internal/usecase"

	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Provide(
		func(cfg *config.Config) (usecase.BackendClient, usecase.CatalogClient) {
			if cfg.MockBackend {
				m := NewMockClient()
				return m, m
			}
			c := NewClient(cfg.BackendURL)
			return c, c
		},
	),
)
