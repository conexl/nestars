package keyboards

import (
	"github.com/mymmrac/telego"

	"nestars-bot/internal/common"
	"nestars-bot/internal/config"
	"nestars-bot/internal/i18n"
)

func PaymentMethods(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.PayHeleket(lang),
					CallbackData:      common.CbPayment + ":Heleket",
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5334653091654439979",
				},
				{
					Text:              i18n.PayYookassa(lang),
					CallbackData:      common.CbPayment + ":Yookassa",
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5363972466857252756",
				},
			},
			{
				{
					Text:              i18n.PayCryptoBot(lang),
					CallbackData:      common.CbPayment + ":CryptoPay",
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5361914370068613491",
				},
				{
					Text:              i18n.PayBalance(lang),
					CallbackData:      common.CbPayment + ":Balance",
					Style:             telego.ButtonStylePrimary,
				},
			},
			{
				BackButton(lang),
			},
		},
	}
}

func PaymentReady(lang i18n.Lang, payURL string) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				PayButton(lang, payURL),
			},
			{
				{
					Text:              i18n.Support(lang),
					URL:               config.Cfg.Links.Support,
					Style:             telego.ButtonStyleSuccess,
					IconCustomEmojiID: "5931409969613116639",
				},
			},
			{
				{
					Text:              i18n.Back(lang),
					CallbackData:      common.CbPaymentMain,
					Style:             telego.ButtonStyleDanger,
					IconCustomEmojiID: "5877536313623711363",
				},
			},
		},
	}
}

func PaymentReadyBalance(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				PayButtonBalance(lang),
			},
			{
				{
					Text:              i18n.Support(lang),
					URL:               config.Cfg.Links.Support,
					Style:             telego.ButtonStyleSuccess,
					IconCustomEmojiID: "5931409969613116639",
				},
			},
			{
				{
					Text:              i18n.Back(lang),
					CallbackData:      common.CbPaymentMain,
					Style:             telego.ButtonStyleDanger,
					IconCustomEmojiID: "5877536313623711363",
				},
			},
		},
	}
}
