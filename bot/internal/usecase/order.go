package usecase

import (
	"context"
	"fmt"
	"nestars-bot/internal/domain"
)

type Order struct {
	backend BackendClient
	repo    domain.OrderRepository
}

func NewOrder(backend BackendClient, repo domain.OrderRepository) *Order {
	return &Order{backend: backend, repo: repo}
}

func (uc *Order) Create(ctx context.Context, in domain.CreateOrderInput) (*domain.Order, error) {
	// Если кастомные звезды без ProductID — считаем здесь или на backend
	if in.ProductType == "Stars" && in.Stars > 0 && in.ProductID == "" {
		in.ProductID = "custom"
	}

	order, err := uc.backend.CreateOrder(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("backend: %w", err)
	}

	// Сохраняем в локальную БД для истории (не критично)
	_ = uc.repo.Create(ctx, order)
	return order, nil
}

func (uc *Order) Pay(ctx context.Context, orderID, method string) (string, error) {
	if method == "Balance" {
		if err := uc.repo.ProcessBalancePayment(ctx, orderID); err != nil {
			return "", err
		}
		return "", nil // мгновенная оплата
	}

	url, err := uc.backend.CreatePayment(ctx, domain.CreatePaymentInput{
		OrderID: orderID,
		Method:  method,
	})
	if err != nil {
		return "", err
	}

	_ = uc.repo.UpdatePayment(ctx, orderID, method)
	_ = uc.repo.UpdatePaymentURL(ctx, orderID, url)
	return url, nil
}

func (uc *Order) Withdraw(ctx context.Context, in domain.WithdrawInput) error {
	return uc.backend.WithdrawStars(ctx, in)
}
