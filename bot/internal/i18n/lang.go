package i18n

import "strings"

// LangStorage — минимальный интерфейс для хранения выбранного языка.
type LangStorage interface {
	GetData(chatID int64, key string) (any, bool)
	SetData(chatID int64, key string, value any)
}

type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
)

func Normalize(code string) Lang {
	code = strings.ToLower(strings.TrimSpace(code))
	if strings.HasPrefix(code, "ru") {
		return RU
	}

	return EN
}

func FromStorage(storage LangStorage, chatID int64) Lang {
	raw, _ := storage.GetData(chatID, "lang")
	if lang, ok := raw.(string); ok && lang != "" {
		return Normalize(lang)
	}

	return EN
}

func Remember(storage LangStorage, chatID int64, code string) Lang {
	current := FromStorage(storage, chatID)
	if current != EN {
		return current
	}

	lang := Normalize(code)
	storage.SetData(chatID, "lang", string(lang))
	return lang
}
