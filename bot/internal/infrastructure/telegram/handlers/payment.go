package handlers

import (
	"context"
	"strconv"
	"strings"

	"nestars-bot/internal/i18n"
	"nestars-bot/internal/keyboards"
	"nestars-bot/internal/usecase"

	"github.com/mymmrac/telego"
)

type Payment struct {
	order  *usecase.Order
	render Render
	fsm    FSM
}

func NewPayment(order *usecase.Order, render Render, fsm FSM) *Payment {
	return &Payment{order: order, render: render, fsm: fsm}
}

func (h *Payment) Handle(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	lang := h.lang(chatID)

	method := strings.TrimPrefix(cb.Data, "payment:")
	orderID, _ := h.fsm.GetData(chatID, "pending_order_id")
	oid, _ := orderID.(string)

	url, err := h.order.Pay(ctx, oid, method)
	if err != nil {
		h.render.Answer(ctx, cb.ID, err.Error(), true)
		return
	}

	amount, _ := h.fsm.GetData(chatID, "pending_amount")
	number, _ := h.fsm.GetData(chatID, "pending_order_number")

	a, _ := amount.(float64)
	n, _ := number.(int)

	text := i18n.PaymentReady(lang, formatAmount(a, "RUB"), strconv.Itoa(n))

	if method == "Balance" {
		h.render.EditMenu(ctx, chatID, text, keyboards.PaymentReadyBalance(lang))
	} else {
		h.render.EditMenu(ctx, chatID, text, keyboards.PaymentReady(lang, url))
	}

	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Payment) lang(chatID int64) i18n.Lang {
	raw, _ := h.fsm.GetData(chatID, "lang")
	if s, ok := raw.(string); ok && s != "" {
		return i18n.Normalize(s)
	}
	return i18n.EN
}
