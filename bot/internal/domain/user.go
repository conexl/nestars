package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	TelegramID   int64     `json:"telegram_id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	LanguageCode string    `json:"language_code"`
	ReferralCode string    `json:"referral_code"`
	ReferredBy   *string   `json:"referred_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type Profile struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"user_id"`
	TelegramID         int64   `json:"telegram_id"`
	ReferralCode       string  `json:"referral_code"`
	BalanceUSD         float64 `json:"balance_usd"`
	BalanceTON         float64 `json:"balance_ton"`
	BoughtStars        int     `json:"bought_stars"`
	BonusStars         int     `json:"bonus_stars"`
	ReferralBonusStars int     `json:"referral_bonus_stars"`
	ReferralCount      int     `json:"referral_count"`
}

type TelegramUser struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	LanguageCode string `json:"language_code"`
}
