package keyboards

import (
	"github.com/mymmrac/telego"

	"nestars-bot/internal/i18n"
)

func WithdrawWhoMenu(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.Self(lang),
					CallbackData:      "withdraw_self",
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "6145278957039131275",
				},
			},
			{
				BackButton(lang),
			},
		},
	}
}

func WithdrawFlowMenu(lang i18n.Lang, recipient string, amount string, ready bool) *telego.InlineKeyboardMarkup {
	buttons := [][]telego.InlineKeyboardButton{
		{
			{Text: withdrawRecipientLabel(lang, recipient), CallbackData: "withdraw_select_user", Style: telego.ButtonStylePrimary, IconCustomEmojiID: "6145278957039131275"},
		},
		{
			{Text: i18n.EnterAmount(lang), CallbackData: "withdraw_amount", Style: telego.ButtonStylePrimary, IconCustomEmojiID: "5879841310902324730"},
		},
	}

	if ready {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			BackButton(lang),
			{Text: i18n.Withdraw(lang), CallbackData: "withdraw_submit", Style: telego.ButtonStyleSuccess},
		})
	} else {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			BackButton(lang),
		})
	}

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func withdrawRecipientLabel(lang i18n.Lang, recipient string) string {
	return i18n.RecipientLabel(lang, recipient)
}
