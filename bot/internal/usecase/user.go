package usecase

import (
	"context"
	"nestars-bot/internal/domain"
)

type User struct {
	repo domain.UserRepository
}

func NewUser(repo domain.UserRepository) *User {
	return &User{repo: repo}
}

func (u *User) CreateOrGet(ctx context.Context, tg *domain.TelegramUser, ref string) (*domain.User, error) {
	return u.repo.CreateOrGet(ctx, tg, ref)
}

func (u *User) Profile(ctx context.Context, telegramID int64) (*domain.User, *domain.Profile, error) {
	return u.repo.GetUserWithProfile(ctx, telegramID)
}
