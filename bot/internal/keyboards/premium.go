package keyboards

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"nestars-bot/internal/common"
	"nestars-bot/internal/i18n"
)

type PremiumPlanItem struct {
	ID       string  `json:"ID"`
	Name     string  `json:"Name"`
	Value    int     `json:"Value"`
	Amount   float64 `json:"Amount"`
	Currency string  `json:"Currency"`
}

type PremiumPlanCatalogResponse struct {
	Data    []PremiumPlanItem `json:"Data"`
	Success bool              `json:"Success"`
}

func PremiumWhoMenu(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.Self(lang),
					CallbackData:      common.CbPremiumSelf,
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

func PremiumFlowMenu(lang i18n.Lang, storage FlowStorage, chatID int64, items []PremiumPlanItem) *telego.InlineKeyboardMarkup {
	recipientRaw, _ := storage.GetData(chatID, "recipient_username")
	productRaw, _ := storage.GetData(chatID, "premium_product_name")

	recipient := i18n.NotSelected(lang)
	if s, ok := recipientRaw.(string); ok && s != "" {
		recipient = i18n.OptionalHandle(s)
	}

	product := i18n.NotSelectedValue(lang)
	if s, ok := productRaw.(string); ok && s != "" {
		product = s
	}

	var buttons [][]telego.InlineKeyboardButton

	if len(items) == 0 {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(i18n.NoProducts(lang)).WithCallbackData(common.CbBack),
		})
	} else {
		row := make([]telego.InlineKeyboardButton, 0, 4)
		for _, item := range items {
			btn := tu.InlineKeyboardButton(PremiumButtonLabelLang(lang, item)).
				WithCallbackData("buy_premium:" + item.ID).
				WithStyle(telego.ButtonStylePrimary).
				WithIconCustomEmojiID("5888993774540951956")
			row = append(row, btn)

			if len(row) == 4 {
				buttons = append(buttons, row)
				row = make([]telego.InlineKeyboardButton, 0, 4)
			}
		}

		if len(row) > 0 {
			buttons = append(buttons, row)
		}
	}

	buttons = append(buttons, []telego.InlineKeyboardButton{
		{
			Text:              i18n.Self(lang),
			CallbackData:      common.CbPremiumSelf,
			Style:             telego.ButtonStylePrimary,
			IconCustomEmojiID: "6145278957039131275",
		},
		{
			Text:              premiumRecipientLabel(lang, recipient),
			CallbackData:      "premium_select_user",
			Style:             telego.ButtonStylePrimary,
			IconCustomEmojiID: "6145278957039131275",
		},
	})

	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: i18n.EnterAmount(lang), CallbackData: common.CbPremiumCustom, Style: telego.ButtonStylePrimary, IconCustomEmojiID: "5879841310902324730"},
	})

	if recipient != i18n.NotSelected(lang) && product != i18n.NotSelectedValue(lang) {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			BackButton(lang),
			{
				Text:              i18n.Continue(lang),
				CallbackData:      "premium_checkout",
				Style:             telego.ButtonStyleSuccess,
				IconCustomEmojiID: "5875506366050734240",
			},
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

func PremiumButtonLabelLang(lang i18n.Lang, item PremiumPlanItem) string {
	if item.Value > 0 {
		return i18n.PremiumLabel(lang, item.Value)
	}
	if item.Name != "" {
		return item.Name
	}
	return item.ID
}

func premiumRecipientLabel(lang i18n.Lang, recipient string) string {
	return i18n.RecipientLabel(lang, recipient)
}
