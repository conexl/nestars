package usecase

import (
	"context"
	"nestars-bot/internal/domain"
)

// Интерфейсы к твоему backend (реализуются в infrastructure/api)
type BackendClient interface {
	CreateOrder(ctx context.Context, in domain.CreateOrderInput) (*domain.Order, error)
	CreatePayment(ctx context.Context, in domain.CreatePaymentInput) (string, error)
	WithdrawStars(ctx context.Context, in domain.WithdrawInput) error
}

type CatalogClient interface {
	GetStarsCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error)
	GetPremiumCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error)
}
