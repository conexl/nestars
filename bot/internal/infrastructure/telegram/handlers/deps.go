package handlers

import (
	"context"

	"github.com/mymmrac/telego"
	"nestars-bot/internal/domain"
)

// Render — абстракция над отрисовкой сообщений бота.
// Реализуется в пакете telegram, что позволяет избежать циклического импорта с handlers.
type Render interface {
	SendMenu(ctx context.Context, chatID int64, text string, kb *telego.InlineKeyboardMarkup)
	EditMenu(ctx context.Context, chatID int64, text string, kb *telego.InlineKeyboardMarkup)
	DeleteMsg(ctx context.Context, chatID int64, msgID int)
	Answer(ctx context.Context, queryID, text string, alert bool)
}

// FSM — абстракция над машиной состояний диалога.
// Реализуется в пакете telegram.
type FSM interface {
	Current(chatID int64) domain.FlowState
	Push(chatID int64, s domain.FlowState)
	Pop(chatID int64) domain.FlowState
	Replace(chatID int64, s domain.FlowState)
	Reset(chatID int64)
	SetData(chatID int64, key string, value any)
	GetData(chatID int64, key string) (any, bool)
	DeleteData(chatID int64, key string)
}
