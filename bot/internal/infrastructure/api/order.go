package api

import (
	"context"
	"nestars-bot/internal/domain"
)

type orderResp struct {
	Data struct {
		ID          string  `json:"ID"`
		Number      int     `json:"Number"`
		ProductType string  `json:"ProductType"`
		Value       int     `json:"Value"`
		FinalAmount float64 `json:"FinalAmount"`
		Currency    string  `json:"Currency"`
		Status      string  `json:"Status"`
	} `json:"Data"`
}

func (c *Client) CreateOrder(ctx context.Context, in domain.CreateOrderInput) (*domain.Order, error) {
	req := struct {
		UserID      string `json:"UserID"`
		Username    string `json:"Username"`
		ProductType string `json:"ProductType"`
		ProductID   string `json:"ProductID,omitempty"`
		Stars       int    `json:"Stars,omitempty"`
		Currency    string `json:"Currency"`
	}{in.UserID, in.Username, in.ProductType, in.ProductID, in.Stars, in.Currency}

	var resp orderResp
	if err := c.post(ctx, "/api/orders", req, &resp); err != nil {
		return nil, err
	}

	o := &domain.Order{
		ID:                resp.Data.ID,
		OrderNumber:       resp.Data.Number,
		UserID:            in.UserID,
		RecipientUsername: in.Username,
		ProductType:       resp.Data.ProductType,
		ProductID:         in.ProductID,
		Amount:            resp.Data.FinalAmount,
		Currency:          resp.Data.Currency,
		Status:            resp.Data.Status,
	}
	if in.Stars > 0 {
		o.StarsAmount = &in.Stars
	}
	return o, nil
}

type paymentResp struct {
	Data struct {
		PaymentURL string `json:"PaymentURL"`
		URL        string `json:"URL"`
		Link       string `json:"Link"`
	} `json:"Data"`
}

func (c *Client) CreatePayment(ctx context.Context, in domain.CreatePaymentInput) (string, error) {
	req := struct {
		OrderID string `json:"OrderID"`
		Method  string `json:"Method"`
	}{in.OrderID, in.Method}

	var resp paymentResp
	if err := c.post(ctx, "/api/payments", req, &resp); err != nil {
		return "", err
	}
	if resp.Data.PaymentURL != "" {
		return resp.Data.PaymentURL, nil
	}
	if resp.Data.URL != "" {
		return resp.Data.URL, nil
	}
	return resp.Data.Link, nil
}

type withdrawReq struct {
	UserID   string `json:"UserID"`
	Username string `json:"Username"`
	Stars    int    `json:"Stars"`
}

func (c *Client) WithdrawStars(ctx context.Context, in domain.WithdrawInput) error {
	req := withdrawReq{UserID: in.UserID, Username: in.Username, Stars: in.Stars}
	return c.post(ctx, "/api/withdraw", req, nil)
}
