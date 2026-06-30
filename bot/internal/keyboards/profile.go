package keyboards

import (
	"github.com/mymmrac/telego"

	"nestars-bot/internal/common"
	"nestars-bot/internal/i18n"
)

func ProfileMenu(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.ReferralProgram(lang),
					CallbackData:      common.CbReferralProgram,
					IconCustomEmojiID: "5409048419211682843",
					Style:             telego.ButtonStyleSuccess,
				},
			},
			{
				BackButton(lang),
			},
		},
	}
}

func ReferralMenu(lang i18n.Lang, link string) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.WithdrawStars(lang),
					CallbackData:      common.CbWithdrawStars,
					IconCustomEmojiID: "5960672896060756972",
					Style:             telego.ButtonStylePrimary,
				},
			},
			{
				BackButton(lang),
			},
		},
	}
}
