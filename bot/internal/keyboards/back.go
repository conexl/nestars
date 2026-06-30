package keyboards

import (
	"github.com/mymmrac/telego"

	"nestars-bot/internal/common"
	"nestars-bot/internal/i18n"
)

func BackButton(lang i18n.Lang) telego.InlineKeyboardButton {
	return telego.InlineKeyboardButton{
		Text:              i18n.Back(lang),
		CallbackData:      common.CbBack,
		Style:             telego.ButtonStyleDanger,
		IconCustomEmojiID: "5877536313623711363",
	}
}

func BackOnly(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				BackButton(lang),
			},
		},
	}
}

func Empty() *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{},
	}
}
