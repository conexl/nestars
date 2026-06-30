package handlers

import (
	"context"
	"strings"

	"nestars-bot/internal/domain"
	"nestars-bot/internal/i18n"
	"nestars-bot/internal/keyboards"
	"nestars-bot/internal/usecase"

	"github.com/mymmrac/telego"
)

type Menu struct {
	userUC *usecase.User
	render Render
	fsm    FSM
}

func NewMenu(userUC *usecase.User, render Render, fsm FSM) *Menu {
	return &Menu{userUC: userUC, render: render, fsm: fsm}
}

func (h *Menu) Start(ctx context.Context, msg *telego.Message) {
	chatID := msg.Chat.ID
	h.fsm.Reset(chatID)

	lang := i18n.Normalize(msg.From.LanguageCode)
	h.fsm.SetData(chatID, "lang", string(lang))

	tgUser := &domain.TelegramUser{
		ID:           msg.From.ID,
		Username:     msg.From.Username,
		FirstName:    msg.From.FirstName,
		LastName:     msg.From.LastName,
		LanguageCode: msg.From.LanguageCode,
	}

	user, err := h.userUC.CreateOrGet(ctx, tgUser, parseRef(msg.Text))
	if err != nil {
		h.render.EditMenu(ctx, chatID, i18n.CreateUserError(lang), keyboards.BackOnly(lang))
		return
	}

	h.fsm.SetData(chatID, "user_uuid", user.ID)
	h.fsm.SetData(chatID, "telegram_id", msg.From.ID)
	h.fsm.SetData(chatID, "username", msg.From.Username)

	h.render.SendMenu(ctx, chatID, i18n.MainMenu(lang), keyboards.MainMenu(lang))
}

func (h *Menu) ShowMain(ctx context.Context, chatID int64) {
	lang := h.lang(chatID)
	h.render.EditMenu(ctx, chatID, i18n.MainMenu(lang), keyboards.MainMenu(lang))
}

func (h *Menu) lang(chatID int64) i18n.Lang {
	raw, _ := h.fsm.GetData(chatID, "lang")
	if s, ok := raw.(string); ok && s != "" {
		return i18n.Normalize(s)
	}
	return i18n.EN
}

func parseRef(text string) string {
	// /start REFCODE
	fields := strings.Fields(text)
	if len(fields) == 2 {
		return fields[1]
	}
	return ""
}
