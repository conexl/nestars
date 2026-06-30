package usecase

import (
	"context"
	"nestars-bot/internal/domain"
)

type Catalog struct {
	client CatalogClient
}

func NewCatalog(client CatalogClient) *Catalog {
	return &Catalog{client: client}
}

func (c *Catalog) Stars(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	return c.client.GetStarsCatalog(ctx, currency)
}

func (c *Catalog) Premium(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	return c.client.GetPremiumCatalog(ctx, currency)
}
