package miniapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const ctxInitDataKey contextKey = "initData"

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.DisableMiniAppAuth {
			ctx := context.WithValue(r.Context(), ctxInitDataKey, &InitData{
				User: &TelegramUser{ID: 1, FirstName: "Demo", Username: "demo", LanguageCode: "ru"},
			})
			next(w, r.WithContext(ctx))
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "tma" {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		data, err := validateInitData(parts[1], s.cfg.Telegram.Token)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}

		if data.User == nil {
			http.Error(w, `{"error":"user not found in init data"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxInitDataKey, data)
		next(w, r.WithContext(ctx))
	}
}

func getInitData(r *http.Request) *InitData {
	data, _ := r.Context().Value(ctxInitDataKey).(*InitData)
	return data
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
