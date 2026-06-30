package handlers

import (
	"context"

	"nestars-bot/internal/i18n"
	"nestars-bot/internal/keyboards"
	"nestars-bot/internal/usecase"
)

type Profile struct {
	userUC *usecase.User
	render Render
	fsm    FSM
}

func NewProfile(userUC *usecase.User, render Render, fsm FSM) *Profile {
	return &Profile{userUC: userUC, render: render, fsm: fsm}
}

func (h *Profile) Show(ctx context.Context, chatID, telegramID int64) {
	lang := h.lang(chatID)

	if telegramID == 0 {
		tidRaw, _ := h.fsm.GetData(chatID, "telegram_id")
		telegramID, _ = tidRaw.(int64)
	}
	if telegramID == 0 {
		h.render.EditMenu(ctx, chatID, i18n.ProfileLoadError(lang), keyboards.BackOnly(lang))
		return
	}

	// Запоминаем ID, чтобы после рестарта не терять профиль.
	h.fsm.SetData(chatID, "telegram_id", telegramID)

	user, prof, err := h.userUC.Profile(ctx, telegramID)
	if err != nil {
		h.render.EditMenu(ctx, chatID, i18n.ProfileLoadError(lang), keyboards.BackOnly(lang))
		return
	}

	createdAt := user.CreatedAt.Format("02.01.2006")
	text := i18n.ProfileText(lang, telegramID, createdAt, prof.BalanceUSD, prof.BalanceTON, prof.BoughtStars)
	h.render.EditMenu(ctx, chatID, text, keyboards.ProfileMenu(lang))
}

func (h *Profile) lang(chatID int64) i18n.Lang {
	raw, _ := h.fsm.GetData(chatID, "lang")
	if s, ok := raw.(string); ok && s != "" {
		return i18n.Normalize(s)
	}
	return i18n.EN
}
