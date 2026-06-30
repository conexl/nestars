package miniapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"nestars-bot/internal/domain"
)

// TelegramUser mirrors the user object inside initData.
type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
	IsPremium    bool   `json:"is_premium"`
}

// InitData contains parsed Telegram WebApp init data.
type InitData struct {
	Raw      string
	QueryID  string
	User     *TelegramUser
	AuthDate time.Time
	Hash     string
}

// parseInitData parses the decoded init data values.
func parseInitData(raw string) (*InitData, error) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil, fmt.Errorf("parse init data: %w", err)
	}

	data := &InitData{Raw: raw}
	data.Hash = values.Get("hash")
	if data.Hash == "" {
		return nil, fmt.Errorf("hash missing")
	}
	data.QueryID = values.Get("query_id")

	if userRaw := values.Get("user"); userRaw != "" {
		var user TelegramUser
		if err := json.Unmarshal([]byte(userRaw), &user); err != nil {
			return nil, fmt.Errorf("parse user: %w", err)
		}
		data.User = &user
	}

	if authDateStr := values.Get("auth_date"); authDateStr != "" {
		var ts int64
		if _, err := fmt.Sscanf(authDateStr, "%d", &ts); err == nil {
			data.AuthDate = time.Unix(ts, 0)
		}
	}

	return data, nil
}

// splitRawQuery splits the raw query string without URL-decoding values.
func splitRawQuery(raw string) (map[string]string, string, error) {
	pairs := make(map[string]string)
	var hash string

	for _, part := range strings.Split(raw, "&") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]
		if key == "hash" {
			hash = value
			continue
		}
		pairs[key] = value
	}

	if hash == "" {
		return nil, "", fmt.Errorf("hash missing")
	}

	return pairs, hash, nil
}

// validateInitData checks the HMAC-SHA256 signature of init data.
func validateInitData(raw string, botToken string) (*InitData, error) {
	data, err := parseInitData(raw)
	if err != nil {
		return nil, err
	}

	rawPairs, hash, err := splitRawQuery(raw)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(rawPairs))
	for k := range rawPairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	checkParts := make([]string, 0, len(keys))
	for _, k := range keys {
		checkParts = append(checkParts, k+"="+rawPairs[k])
	}
	checkString := strings.Join(checkParts, "\n")

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(botToken))

	h := hmac.New(sha256.New, secret.Sum(nil))
	h.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(hash), []byte(expectedHash)) {
		return nil, fmt.Errorf("invalid init data signature")
	}

	// Reject data older than 24 hours to prevent replay attacks.
	if !data.AuthDate.IsZero() && time.Since(data.AuthDate) > 24*time.Hour {
		return nil, fmt.Errorf("init data expired")
	}

	return data, nil
}

// toDomainUser converts parsed TelegramUser to domain.TelegramUser.
func (u *TelegramUser) toDomainUser() *domain.TelegramUser {
	return &domain.TelegramUser{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Username:     u.Username,
		LanguageCode: u.LanguageCode,
	}
}
