package keyboards

import (
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"nestars-bot/internal/common"
	"nestars-bot/internal/i18n"
)

type StarsItem struct {
	ID       string  `json:"ID"`
	Name     string  `json:"Name"`
	Value    int     `json:"Value"`
	Amount   float64 `json:"Amount"`
	Currency string  `json:"Currency"`
}

type StarsCatalogResponse struct {
	Data    []StarsItem `json:"Data"`
	Success bool        `json:"Success"`
}

// FlowStorage позволяет клавиатурам читать состояние без прямой зависимости от FSM.
type FlowStorage interface {
	GetData(chatID int64, key string) (any, bool)
}

func StarsWhoMenu(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.Self(lang),
					CallbackData:      common.CbStarsSelf,
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

func StarsFlowMenu(lang i18n.Lang, storage FlowStorage, chatID int64, items []StarsItem) *telego.InlineKeyboardMarkup {
	recipientRaw, _ := storage.GetData(chatID, "recipient_username")
	productRaw, _ := storage.GetData(chatID, "stars_product_name")

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
			btn := tu.InlineKeyboardButton(starsButtonLabel(item)).
				WithCallbackData("buy_stars:" + item.ID).
				WithStyle(telego.ButtonStyleSuccess).
				WithIconCustomEmojiID("5954135079662916434")
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
		{Text: starsRecipientLabel(lang, recipient), CallbackData: "stars_select_user", Style: telego.ButtonStylePrimary, IconCustomEmojiID: "6145278957039131275"},
	})

	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: i18n.EnterAmount(lang), CallbackData: common.CbStarsCustom, Style: telego.ButtonStylePrimary, IconCustomEmojiID: "5879841310902324730"},
	})

	if recipient != i18n.NotSelected(lang) && product != i18n.NotSelectedValue(lang) {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			BackButton(lang),
			{
				Text:              i18n.Continue(lang),
				CallbackData:      "stars_checkout",
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

func starsButtonLabel(item StarsItem) string {
	if item.Value > 0 {
		return strconv.Itoa(item.Value)
	}

	if item.Name != "" {
		return item.Name
	}

	return item.ID
}

func formatPrice(amount float64, currency string) string {
	switch currency {
	case "RUB":
		return fmt.Sprintf("%2.f₽", amount)
	case "USD":
		return fmt.Sprintf("$%2.f", amount)
	case "EUR":
		return fmt.Sprintf("€%2.f", amount)
	default:
		return fmt.Sprintf("%f %s", amount, currency)
	}
}

func starsRecipientLabel(lang i18n.Lang, recipient string) string {
	return i18n.RecipientLabel(lang, recipient)
}
