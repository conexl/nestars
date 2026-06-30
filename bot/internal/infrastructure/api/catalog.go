package api

import (
	"context"
	"nestars-bot/internal/domain"
)

type catalogResp struct {
	Data    []domain.CatalogItem `json:"Data"`
	Success bool                 `json:"Success"`
}

func (c *Client) GetStarsCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	var r catalogResp
	if err := c.get(ctx, "/api/catalog/stars?currency="+currency, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) GetPremiumCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	var r catalogResp
	if err := c.get(ctx, "/api/catalog/premium?currency="+currency, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}
