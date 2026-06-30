package domain

import "context"

type UserRepository interface {
	CreateOrGet(ctx context.Context, user *TelegramUser, refCode string) (*User, error)
	GetUserWithProfile(ctx context.Context, telegramID int64) (*User, *Profile, error)
	GetProfile(ctx context.Context, userID string) (*Profile, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	UpdatePayment(ctx context.Context, orderID, method string) error
	UpdatePaymentURL(ctx context.Context, orderID, url string) error
	ProcessBalancePayment(ctx context.Context, orderID string) error
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
}
