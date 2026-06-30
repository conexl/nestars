package keyboards

import (
	"github.com/mymmrac/telego"

	"nestars-bot/internal/common"
	"nestars-bot/internal/config"
	"nestars-bot/internal/i18n"
)

func MainMenu(lang i18n.Lang) *telego.InlineKeyboardMarkup {
	kb := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{
					Text:              i18n.Stars(lang),
					CallbackData:      common.CbBuyStars,
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5954135079662916434",
				},
				{
					Text:              i18n.Premium(lang),
					CallbackData:      common.CbBuyPremium,
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5888993774540951956",
				},
			},
			{
				{
					Text:              i18n.Profile(lang),
					CallbackData:      common.CbProfile,
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "6145278957039131275",
				},
				{
					Text:              i18n.Giveaways(lang),
					URL:               config.Cfg.Links.NewsChannel,
					Style:             telego.ButtonStylePrimary,
					IconCustomEmojiID: "5879585266426973039",
				},
			},
			{
				{
					Text:              i18n.Support(lang),
					URL:               config.Cfg.Links.Support,
					Style:             telego.ButtonStyleSuccess,
					IconCustomEmojiID: "5931409969613116639",
				},
			},
		},
	}

	if config.Cfg.MiniAppURL != "" {
		kb.InlineKeyboard = append([][]telego.InlineKeyboardButton{{
			{
				Text:   i18n.OpenShop(lang),
				WebApp: &telego.WebAppInfo{URL: config.Cfg.MiniAppURL},
				Style:  telego.ButtonStylePrimary,
			},
		}}, kb.InlineKeyboard...)
	}

	return kb
}
