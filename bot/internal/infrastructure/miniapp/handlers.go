package miniapp

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"nestars-bot/internal/domain"
)

type authRequest struct {
	InitData string `json:"init_data"`
}

type authResponse struct {
	User    *domain.User    `json:"user"`
	Profile *domain.Profile `json:"profile"`
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	if s.cfg.DisableMiniAppAuth {
		user, profile, err := s.getOrCreateUser(r.Context(), &TelegramUser{ID: 1, FirstName: "Demo", Username: "demo", LanguageCode: "ru"})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, authResponse{User: user, Profile: profile})
		return
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	data, err := validateInitData(req.InitData, s.cfg.Telegram.Token)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	user, profile, err := s.getOrCreateUser(r.Context(), data.User)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: user, Profile: profile})
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	data := getInitData(r)
	if data == nil || data.User == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, profile, err := s.getOrCreateUser(r.Context(), data.User)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: user, Profile: profile})
}

func (s *Server) handleCatalogStars(w http.ResponseWriter, r *http.Request) {
	items, err := s.catalog.Stars(r.Context(), "RUB")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCatalogPremium(w http.ResponseWriter, r *http.Request) {
	items, err := s.catalog.Premium(r.Context(), "RUB")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type createOrderRequest struct {
	Username    string `json:"username"`
	ProductType string `json:"product_type"`
	ProductID   string `json:"product_id"`
	Stars       int    `json:"stars"`
	Currency    string `json:"currency"`
}

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	data := getInitData(r)
	if data == nil || data.User == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	user, _, err := s.getOrCreateUser(r.Context(), data.User)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	order, err := s.order.Create(r.Context(), domain.CreateOrderInput{
		UserID:      user.ID,
		Username:    strings.TrimPrefix(req.Username, "@"),
		ProductType: req.ProductType,
		ProductID:   req.ProductID,
		Stars:       req.Stars,
		Currency:    req.Currency,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

type payOrderRequest struct {
	Method string `json:"method"`
}

type payOrderResponse struct {
	PaymentURL string `json:"payment_url,omitempty"`
}

func (s *Server) handlePayOrder(w http.ResponseWriter, r *http.Request) {
	data := getInitData(r)
	if data == nil || data.User == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orderID := strings.TrimPrefix(r.URL.Path, "/api/miniapp/orders/")
	orderID = strings.TrimSuffix(orderID, "/pay")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order id required"})
		return
	}

	var req payOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	url, err := s.order.Pay(r.Context(), orderID, req.Method)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, payOrderResponse{PaymentURL: url})
}

func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	data := getInitData(r)
	if data == nil || data.User == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, _, err := s.getOrCreateUser(r.Context(), data.User)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	orders, err := s.orderRepo.FindByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

type withdrawRequest struct {
	Username string `json:"username"`
	Stars    int    `json:"stars"`
}

func (s *Server) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	data := getInitData(r)
	if data == nil || data.User == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req withdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	user, _, err := s.getOrCreateUser(r.Context(), data.User)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := s.order.Withdraw(r.Context(), domain.WithdrawInput{
		UserID:   user.ID,
		Username: strings.TrimPrefix(req.Username, "@"),
		Stars:    req.Stars,
	}); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (s *Server) getOrCreateUser(ctx context.Context, tgUser *TelegramUser) (*domain.User, *domain.Profile, error) {
	user, err := s.userUC.CreateOrGet(ctx, tgUser.toDomainUser(), "")
	if err != nil {
		return nil, nil, err
	}
	_, profile, err := s.userUC.Profile(ctx, user.TelegramID)
	if err != nil {
		return user, nil, err
	}
	return user, profile, nil
}
