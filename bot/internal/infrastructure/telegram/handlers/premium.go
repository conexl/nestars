package handlers

import (
	"context"
	"strings"

	"github.com/mymmrac/telego"

	"nestars-bot/internal/domain"
	"nestars-bot/internal/i18n"
	"nestars-bot/internal/keyboards"
	"nestars-bot/internal/usecase"
)

type Premium struct {
	catalog *usecase.Catalog
	order   *usecase.Order
	render  Render
	fsm     FSM
}

func NewPremium(catalog *usecase.Catalog, order *usecase.Order, render Render, fsm FSM) *Premium {
	return &Premium{catalog: catalog, order: order, render: render, fsm: fsm}
}

func (h *Premium) ShowFlow(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	items, _ := h.catalog.Premium(ctx, "RUB")

	prodRaw, _ := h.fsm.GetData(chatID, "premium_product_name")
	recipientRaw, _ := h.fsm.GetData(chatID, "recipient_username")

	prod := i18n.NotSelectedValue(lang)
	if p, ok := prodRaw.(string); ok && p != "" {
		prod = p
	}
	rec := i18n.NotSelected(lang)
	if r, ok := recipientRaw.(string); ok && r != "" {
		rec = i18n.OptionalHandle(r)
	}

	text := i18n.PremiumFlow(lang, rec, prod)
	h.render.EditMenu(ctx, chatID, text, keyboards.PremiumFlowMenu(lang, h.fsm, chatID, adaptPremiumItems(items)))
}

func (h *Premium) SelectSelf(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	if cb.From.Username == "" {
		h.render.Answer(ctx, cb.ID, i18n.UsernameRequired(h.lang(chatID)), true)
		return
	}
	h.fsm.SetData(chatID, "recipient_username", cb.From.Username)
	h.ShowFlow(ctx, chatID)
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Premium) ShowRecipient(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	h.render.EditMenu(ctx, chatID, i18n.PremiumWho(lang), keyboards.BackOnly(lang))
}

func (h *Premium) HandleRecipient(ctx context.Context, msg *telego.Message) {
	chatID := msg.Chat.ID
	lang := h.lang(chatID)

	username := parseUsername(msg.Text)
	if username == "" {
		h.render.EditMenu(ctx, chatID, i18n.InvalidUsername(lang), keyboards.BackOnly(lang))
		return
	}

	h.fsm.SetData(chatID, "recipient_username", username)
	h.fsm.Pop(chatID)
	h.ShowFlow(ctx, chatID)
}

func (h *Premium) SelectProduct(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	id := strings.TrimPrefix(cb.Data, "buy_premium:")

	items, _ := h.catalog.Premium(ctx, "RUB")
	lang := h.lang(chatID)
	var label string
	for _, it := range items {
		if it.ID == id {
			label = i18n.PremiumLabel(lang, it.Value) + " " + i18n.MonthWord(lang, it.Value)
			h.fsm.SetData(chatID, "premium_months", it.Value)
			break
		}
	}
	if label == "" {
		label = id
	}

	h.fsm.SetData(chatID, "premium_product_id", id)
	h.fsm.SetData(chatID, "premium_product_name", label)
	h.ShowFlow(ctx, chatID)
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Premium) Checkout(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	lang := h.lang(chatID)

	userUUID, _ := h.fsm.GetData(chatID, "user_uuid")
	recipient, _ := h.fsm.GetData(chatID, "recipient_username")
	prodID, _ := h.fsm.GetData(chatID, "premium_product_id")

	uuid, _ := userUUID.(string)
	rec, _ := recipient.(string)
	pid, _ := prodID.(string)

	if rec == "" {
		h.render.Answer(ctx, cb.ID, i18n.SelectRecipientAndPackage(lang), true)
		return
	}
	if pid == "" {
		h.render.Answer(ctx, cb.ID, i18n.SelectRecipientAndPackage(lang), true)
		return
	}

	in := domain.CreateOrderInput{
		UserID:      uuid,
		Username:    rec,
		ProductType: "Premium",
		ProductID:   pid,
		Currency:    "RUB",
	}

	order, err := h.order.Create(ctx, in)
	if err != nil {
		h.render.Answer(ctx, cb.ID, err.Error(), true)
		return
	}

	h.fsm.SetData(chatID, "pending_order_id", order.ID)
	h.fsm.SetData(chatID, "pending_order_number", order.OrderNumber)
	h.fsm.SetData(chatID, "pending_amount", order.Amount)
	h.fsm.Replace(chatID, domain.StatePremiumPayment)

	productName := h.productName(lang, chatID, order)
	text := i18n.PaymentScreen(lang, true, "@"+rec, productName, formatAmount(order.Amount, order.Currency))
	h.render.EditMenu(ctx, chatID, text, keyboards.PaymentMethods(lang))
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Premium) productName(lang i18n.Lang, chatID int64, o *domain.Order) string {
	if raw, _ := h.fsm.GetData(chatID, "premium_product_name"); raw != nil {
		if s, ok := raw.(string); ok && s != "" {
			return s
		}
	}
	return o.ProductID
}

func (h *Premium) lang(chatID int64) i18n.Lang {
	raw, _ := h.fsm.GetData(chatID, "lang")
	if s, ok := raw.(string); ok && s != "" {
		return i18n.Normalize(s)
	}
	return i18n.EN
}

func adaptPremiumItems(items []domain.CatalogItem) []keyboards.PremiumPlanItem {
	out := make([]keyboards.PremiumPlanItem, len(items))
	for i, it := range items {
		out[i] = keyboards.PremiumPlanItem{
			ID:       it.ID,
			Name:     it.Name,
			Value:    it.Value,
			Amount:   it.Amount,
			Currency: it.Currency,
		}
	}
	return out
}

func parseUsername(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "https://t.me/") {
		return strings.TrimPrefix(text, "https://t.me/")
	}
	if strings.HasPrefix(text, "@") {
		return strings.TrimPrefix(text, "@")
	}
	return ""
}
