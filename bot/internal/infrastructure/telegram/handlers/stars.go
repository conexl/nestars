package handlers

import (
	"context"
	"strconv"
	"strings"

	"nestars-bot/internal/domain"
	"nestars-bot/internal/i18n"
	"nestars-bot/internal/keyboards"
	"nestars-bot/internal/usecase"

	"github.com/mymmrac/telego"
)

type Stars struct {
	catalog *usecase.Catalog
	order   *usecase.Order
	render  Render
	fsm     FSM
}

func NewStars(catalog *usecase.Catalog, order *usecase.Order, render Render, fsm FSM) *Stars {
	return &Stars{catalog: catalog, order: order, render: render, fsm: fsm}
}

func (h *Stars) ShowFlow(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	items, _ := h.catalog.Stars(ctx, "RUB")

	// Восстанавливаем выбор пользователя
	prodRaw, _ := h.fsm.GetData(chatID, "stars_product_name")
	recipientRaw, _ := h.fsm.GetData(chatID, "recipient_username")

	prod := i18n.NotSelectedValue(lang)
	if p, ok := prodRaw.(string); ok && p != "" {
		prod = p
	}
	rec := i18n.NotSelected(lang)
	if r, ok := recipientRaw.(string); ok && r != "" {
		rec = i18n.OptionalHandle(r)
	}

	text := i18n.StarsFlow(lang, rec, prod)
	h.render.EditMenu(ctx, chatID, text, keyboards.StarsFlowMenu(lang, h.fsm, chatID, adaptItems(items)))
}

func (h *Stars) SelectSelf(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	if cb.From.Username == "" {
		h.render.Answer(ctx, cb.ID, i18n.UsernameRequired(h.lang(chatID)), true)
		return
	}
	h.fsm.SetData(chatID, "recipient_username", cb.From.Username)
	h.ShowFlow(ctx, chatID)
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Stars) SelectProduct(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	id := strings.TrimPrefix(cb.Data, "buy_stars:")
	h.fsm.SetData(chatID, "stars_product_id", id)
	h.fsm.SetData(chatID, "stars_product_name", id) // или маппинг на имя
	h.fsm.SetData(chatID, "custom_stars", 0)
	h.ShowFlow(ctx, chatID)
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Stars) ShowCustomAmount(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	h.render.EditMenu(ctx, chatID, i18n.StarsCustomAmount(lang), keyboards.BackOnly(lang))
}

func (h *Stars) HandleCustomAmount(ctx context.Context, msg *telego.Message) {
	chatID := msg.Chat.ID
	lang := h.lang(chatID)

	amount, err := strconv.Atoi(strings.TrimSpace(msg.Text))
	if err != nil || amount < 50 || amount > 1_000_000 {
		h.render.EditMenu(ctx, chatID, i18n.InvalidStarsAmount(lang), keyboards.BackOnly(lang))
		return
	}

	h.fsm.SetData(chatID, "custom_stars", amount)
	h.fsm.SetData(chatID, "stars_product_id", "")
	h.fsm.SetData(chatID, "stars_product_name", strconv.Itoa(amount))
	h.fsm.Pop(chatID) // выходим из состояния ввода

	h.ShowFlow(ctx, chatID)
}

func (h *Stars) ShowRecipient(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	h.render.EditMenu(ctx, chatID, i18n.EnterFriendUsername(lang), keyboards.BackOnly(lang))
}

func (h *Stars) Checkout(ctx context.Context, cb *telego.CallbackQuery) {
	chatID := cb.Message.GetChat().ID
	lang := h.lang(chatID)

	userUUID, _ := h.fsm.GetData(chatID, "user_uuid")
	recipient, _ := h.fsm.GetData(chatID, "recipient_username")
	prodID, _ := h.fsm.GetData(chatID, "stars_product_id")
	custom, _ := h.fsm.GetData(chatID, "custom_stars")

	uuid, _ := userUUID.(string)
	rec, _ := recipient.(string)
	pid, _ := prodID.(string)
	cst, _ := custom.(int)

	if rec == "" {
		h.render.Answer(ctx, cb.ID, i18n.SelectRecipientAndPackage(lang), true)
		return
	}

	in := domain.CreateOrderInput{
		UserID:      uuid,
		Username:    rec,
		ProductType: "Stars",
		ProductID:   pid,
		Stars:       cst,
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
	h.fsm.Replace(chatID, domain.StateStarsPayment)

	text := i18n.PaymentScreen(lang, false, "@"+rec, h.productName(lang, order), formatAmount(order.Amount, order.Currency))
	h.render.EditMenu(ctx, chatID, text, keyboards.PaymentMethods(lang))
	h.render.Answer(ctx, cb.ID, "", false)
}

func (h *Stars) productName(lang i18n.Lang, o *domain.Order) string {
	if o.StarsAmount != nil {
		return i18n.StarsLabel(lang, *o.StarsAmount)
	}
	return o.ProductID
}

func (h *Stars) lang(chatID int64) i18n.Lang {
	raw, _ := h.fsm.GetData(chatID, "lang")
	if s, ok := raw.(string); ok && s != "" {
		return i18n.Normalize(s)
	}
	return i18n.EN
}

func adaptItems(items []domain.CatalogItem) []keyboards.StarsItem {
	// маппинг domain.CatalogItem -> keyboards.StarsItem
	out := make([]keyboards.StarsItem, len(items))
	for i, it := range items {
		out[i] = keyboards.StarsItem{
			ID:       it.ID,
			Name:     it.Name,
			Value:    it.Value,
			Amount:   it.Amount,
			Currency: it.Currency,
		}
	}
	return out
}

func formatAmount(a float64, c string) string {
	if c == "RUB" {
		return strconv.FormatFloat(a, 'f', 2, 64) + " ₽"
	}
	return strconv.FormatFloat(a, 'f', 2, 64) + " " + c
}
