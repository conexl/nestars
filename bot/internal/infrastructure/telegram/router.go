package telegram

import (
	"context"
	"log"
	"strings"

	"nestars-bot/internal/domain"
	"nestars-bot/internal/infrastructure/telegram/handlers"

	"github.com/mymmrac/telego"
)

type Router struct {
	bot     *telego.Bot
	fsm     *FSM
	render  *Render
	menu    *handlers.Menu
	profile *handlers.Profile
	stars   *handlers.Stars
	premium *handlers.Premium
	payment *handlers.Payment
}

func NewRouter(
	bot *telego.Bot,
	fsm *FSM,
	render *Render,
	menu *handlers.Menu,
	profile *handlers.Profile,
	stars *handlers.Stars,
	premium *handlers.Premium,
	payment *handlers.Payment,
) *Router {
	return &Router{
		bot: bot, fsm: fsm, render: render,
		menu: menu, profile: profile,
		stars: stars, premium: premium, payment: payment,
	}
}

func (r *Router) Dispatch(ctx context.Context, up telego.Update) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("router panic: %v", rec)
		}
	}()

	if up.CallbackQuery != nil {
		r.handleCallback(ctx, up.CallbackQuery)
		return
	}
	if up.Message == nil {
		return
	}

	msg := up.Message
	chatID := msg.Chat.ID

	// /start
	if strings.TrimSpace(msg.Text) == "/start" {
		r.menu.Start(ctx, msg)
		return
	}
	if strings.TrimSpace(msg.Text) == "/menu" {
		r.fsm.Reset(chatID)
		r.menu.ShowMain(ctx, chatID)
		return
	}

	// Текстовый ввод по состоянию
	state := r.fsm.Current(chatID)
	switch state {
	case domain.StateStarsCustom:
		r.stars.HandleCustomAmount(ctx, msg)
	case domain.StateStarsWho:
		// TODO: ввод username для Stars
	case domain.StatePremiumWho:
		r.premium.HandleRecipient(ctx, msg)
	case domain.StatePremiumCustom:
		r.premium.HandleCustomAmount(ctx, msg)
	case domain.StateWithdrawWho, domain.StateWithdrawAmount:
		// TODO: withdraw
	}
}

func (r *Router) handleCallback(ctx context.Context, cb *telego.CallbackQuery) {
	if cb.Message == nil {
		return
	}
	chatID := cb.Message.GetChat().ID
	msgID := cb.Message.GetMessageID()
	r.render.lastMsg[chatID] = msgID // обновляем активное сообщение
	data := cb.Data

	switch {
	// Главное меню
	case data == "buy_stars":
		r.fsm.Push(chatID, domain.StateStarsFlow)
		r.stars.ShowFlow(ctx, chatID)
	case data == "buy_premium":
		r.fsm.Push(chatID, domain.StatePremiumFlow)
		r.premium.ShowFlow(ctx, chatID)
	case data == "profile":
		r.fsm.Push(chatID, domain.StateProfile)
		r.profile.Show(ctx, chatID, cb.From.ID)
	case data == "back":
		r.handleBack(ctx, chatID, cb)
	case data == "stars_self":
		r.stars.SelectSelf(ctx, cb)
	case data == "stars_select_user":
		r.fsm.Push(chatID, domain.StateStarsWho)
		r.stars.ShowRecipient(ctx, chatID)
	case data == "stars_custom":
		r.fsm.Push(chatID, domain.StateStarsCustom)
		r.stars.ShowCustomAmount(ctx, chatID)
	case strings.HasPrefix(data, "buy_stars:"):
		r.stars.SelectProduct(ctx, cb)
	case data == "stars_checkout":
		r.stars.Checkout(ctx, cb)
	case data == "premium_self":
		r.premium.SelectSelf(ctx, cb)
	case data == "premium_select_user":
		r.fsm.Push(chatID, domain.StatePremiumWho)
		r.premium.ShowRecipient(ctx, chatID)
	case data == "premium_custom":
		r.fsm.Push(chatID, domain.StatePremiumCustom)
		r.premium.ShowCustomAmount(ctx, chatID)
	case strings.HasPrefix(data, "buy_premium:"):
		r.premium.SelectProduct(ctx, cb)
	case data == "premium_checkout":
		r.premium.Checkout(ctx, cb)
	case strings.HasPrefix(data, "payment:"):
		r.payment.Handle(ctx, cb)
	default:
		r.render.Answer(ctx, cb.ID, "Unknown", false)
	}
}

func (r *Router) handleBack(ctx context.Context, chatID int64, cb *telego.CallbackQuery) {
	current := r.fsm.Current(chatID)
	switch current {
	case domain.StateStarsFlow:
		r.fsm.Reset(chatID)
		r.menu.ShowMain(ctx, chatID)
	case domain.StateStarsWho, domain.StateStarsCustom:
		r.fsm.Pop(chatID)
		r.stars.ShowFlow(ctx, chatID)
	case domain.StateStarsPayment:
		r.fsm.Replace(chatID, domain.StateStarsFlow)
		r.stars.ShowFlow(ctx, chatID)
	case domain.StatePremiumFlow:
		r.fsm.Reset(chatID)
		r.menu.ShowMain(ctx, chatID)
	case domain.StatePremiumWho, domain.StatePremiumCustom:
		r.fsm.Pop(chatID)
		r.premium.ShowFlow(ctx, chatID)
	case domain.StatePremiumPayment:
		r.fsm.Replace(chatID, domain.StatePremiumFlow)
		r.premium.ShowFlow(ctx, chatID)
	case domain.StateProfile:
		r.fsm.Pop(chatID)
		r.menu.ShowMain(ctx, chatID)
	default:
		r.fsm.Pop(chatID)
		r.menu.ShowMain(ctx, chatID)
	}
	r.render.Answer(ctx, cb.ID, "", false)
}
