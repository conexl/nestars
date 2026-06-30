package api

import (
	"context"
	"strconv"
	"strings"

	"nestars-bot/internal/domain"
	"nestars-bot/internal/usecase"
)

// MockClient — заглушка backend, чтобы тестировать UI без реального сервиса.
type MockClient struct{}

var _ usecase.BackendClient = (*MockClient)(nil)
var _ usecase.CatalogClient = (*MockClient)(nil)

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) GetStarsCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	return []domain.CatalogItem{
		{ID: "stars_100", Name: "100 Stars", Value: 100, Amount: 100, Currency: currency},
		{ID: "stars_500", Name: "500 Stars", Value: 500, Amount: 450, Currency: currency},
		{ID: "stars_1000", Name: "1000 Stars", Value: 1000, Amount: 850, Currency: currency},
		{ID: "stars_5000", Name: "5000 Stars", Value: 5000, Amount: 4000, Currency: currency},
	}, nil
}

func (m *MockClient) GetPremiumCatalog(ctx context.Context, currency string) ([]domain.CatalogItem, error) {
	return []domain.CatalogItem{
		{ID: "premium_3", Name: "3 months", Value: 3, Amount: 500, Currency: currency},
		{ID: "premium_6", Name: "6 months", Value: 6, Amount: 900, Currency: currency},
		{ID: "premium_12", Name: "12 months", Value: 12, Amount: 1500, Currency: currency},
	}, nil
}

func (m *MockClient) CreateOrder(ctx context.Context, in domain.CreateOrderInput) (*domain.Order, error) {
	amount := m.estimateAmount(in)

	order := &domain.Order{
		ID:                "order_" + strconv.FormatInt(int64(len(in.Username)+len(in.ProductType)), 36),
		OrderNumber:       1001,
		UserID:            in.UserID,
		RecipientUsername: in.Username,
		ProductType:       in.ProductType,
		ProductID:         in.ProductID,
		Amount:            amount,
		Currency:          in.Currency,
		Status:            "pending",
	}
	if in.ProductType == "Stars" && in.Stars > 0 {
		order.StarsAmount = &in.Stars
	}
	return order, nil
}

func (m *MockClient) CreatePayment(ctx context.Context, in domain.CreatePaymentInput) (string, error) {
	return "https://example.com/pay?order=" + in.OrderID + "&method=" + in.Method, nil
}

func (m *MockClient) WithdrawStars(ctx context.Context, in domain.WithdrawInput) error {
	return nil
}

func (m *MockClient) estimateAmount(in domain.CreateOrderInput) float64 {
	if in.ProductType == "Stars" {
		if in.Stars > 0 {
			return float64(in.Stars) * 1.0
		}
		if strings.HasPrefix(in.ProductID, "stars_") {
			v, _ := strconv.Atoi(strings.TrimPrefix(in.ProductID, "stars_"))
			return float64(v) * 1.0
		}
		return 100
	}
	if in.ProductType == "Premium" {
		switch in.ProductID {
		case "premium_3":
			return 500
		case "premium_6":
			return 900
		case "premium_12":
			return 1500
		}
		if strings.HasPrefix(in.ProductID, "premium_") {
			months, _ := strconv.Atoi(strings.TrimPrefix(in.ProductID, "premium_"))
			return float64(months) * 170
		}
	}
	return 0
}
