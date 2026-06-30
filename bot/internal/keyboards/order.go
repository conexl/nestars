package keyboards

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mymmrac/telego"

	"nestars-bot/internal/common"

	"nestars-bot/internal/config"
	"nestars-bot/internal/i18n"
)

func PayButton(lang i18n.Lang, url string) telego.InlineKeyboardButton {
	return telego.InlineKeyboardButton{
		Text:              i18n.PayOrder(lang),
		URL:               url,
		Style:             telego.ButtonStylePrimary,
		IconCustomEmojiID: "5445353829304387411",
	}
}

func PayButtonBalance(lang i18n.Lang) telego.InlineKeyboardButton {
	return telego.InlineKeyboardButton{
		Text:              i18n.PayOrder(lang),
		CallbackData:      common.CbPaymentBalance,
		Style:             telego.ButtonStylePrimary,
		IconCustomEmojiID: "5445353829304387411",
	}
}

func CancelOrderButton(lang i18n.Lang, orderID string) telego.InlineKeyboardButton {
	return telego.InlineKeyboardButton{
		Text:         i18n.CancelOrder(lang),
		CallbackData: common.CbCancelOrder + ":" + orderID,
	}
}

func ProblemsWithOrderButton(lang i18n.Lang, orderID string) telego.InlineKeyboardButton {
	message := fmt.Sprintf("Hello! I have a problem with order %s", orderID)
	if lang == i18n.RU {
		message = fmt.Sprintf("Здравствуйте! У меня возникла проблема с заказом %s", orderID)
	}
	text := url.QueryEscape(message)
	text = strings.ReplaceAll(text, "+", "%20")

	return telego.InlineKeyboardButton{
		Text: i18n.Support(lang),
		URL:  config.Cfg.Links.Support + "?text=" + text,
	}
}

func Order(lang i18n.Lang, payURL string, orderID string) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				PayButton(lang, payURL),
			},
			{
				CancelOrderButton(lang, orderID),
			},
			{
				ProblemsWithOrderButton(lang, orderID),
			},
		},
	}
}
