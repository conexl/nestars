package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"

	"nestars-bot/internal/config"
)

var Module = fx.Module("persistence",
	fx.Provide(
		NewMongoClient,
		NewDatabase,
		NewUserRepo,
		NewOrderRepo,
	),
)

// NewMongoClient создаёт и подключает клиент MongoDB.
// Жизненный цикл fx обеспечивает корректное закрытие соединения.
func NewMongoClient(lc fx.Lifecycle, cfg *config.Config) (*mongo.Client, error) {
	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is not set")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx, nil); err != nil {
				return fmt.Errorf("ping mongo: %w", err)
			}
			if err := EnsureIndexes(ctx, client.Database(cfg.MongoDB)); err != nil {
				return fmt.Errorf("ensure indexes: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return client, nil
}

// NewDatabase возвращает рабочую базу данных MongoDB.
func NewDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.MongoDB)
}
