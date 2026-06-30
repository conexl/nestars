package telegram

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const menuPhoto = "https://i.ibb.co/4nPw0fcY/L6Zmp-Y8.png"

type Render struct {
	bot     *telego.Bot
	lastMsg map[int64]int
}

func NewRender(bot *telego.Bot) *Render {
	return &Render{bot: bot, lastMsg: make(map[int64]int)}
}

func (r *Render) SendMenu(ctx context.Context, chatID int64, text string, kb *telego.InlineKeyboardMarkup) {
	msg, err := r.bot.SendPhoto(ctx, &telego.SendPhotoParams{
		ChatID:      tu.ID(chatID),
		Photo:       telego.InputFile(tu.FileFromURL(menuPhoto)),
		Caption:     text,
		ParseMode:   telego.ModeHTML,
		ReplyMarkup: kb,
	})
	if err != nil {
		log.Printf("send menu: %v", err)
		return
	}
	r.lastMsg[chatID] = msg.MessageID
}

func (r *Render) EditMenu(ctx context.Context, chatID int64, text string, kb *telego.InlineKeyboardMarkup) {
	msgID := r.lastMsg[chatID]
	if msgID == 0 {
		r.SendMenu(ctx, chatID, text, kb)
		return
	}
	_, err := r.bot.EditMessageCaption(ctx, &telego.EditMessageCaptionParams{
		ChatID:      tu.ID(chatID),
		MessageID:   msgID,
		Caption:     text,
		ParseMode:   telego.ModeHTML,
		ReplyMarkup: kb,
	})
	if err != nil {
		r.SendMenu(ctx, chatID, text, kb)
	}
}

func (r *Render) DeleteMsg(ctx context.Context, chatID int64, msgID int) {
	_ = r.bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: msgID,
	})
}

func (r *Render) Answer(ctx context.Context, queryID, text string, alert bool) {
	_ = r.bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: queryID,
		Text:            text,
		ShowAlert:       alert,
	})
}
