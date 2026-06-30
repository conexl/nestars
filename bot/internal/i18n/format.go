package i18n

import "strconv"

func MonthWord(lang Lang, value int) string {
	if lang == EN {
		if value == 1 {
			return "month"
		}
		return "months"
	}

	mod10 := value % 10
	mod100 := value % 100

	switch {
	case mod10 == 1 && mod100 != 11:
		return "месяц"
	case mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14):
		return "месяца"
	default:
		return "месяцев"
	}
}

func StarsWord(lang Lang, value int) string {
	if lang == EN {
		if value == 1 {
			return "star"
		}
		return "stars"
	}

	mod10 := value % 10
	mod100 := value % 100

	switch {
	case mod10 == 1 && mod100 != 11:
		return "звезда"
	case mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14):
		return "звезды"
	default:
		return "звёзд"
	}
}

func PremiumLabel(lang Lang, value int) string {
	return strconv.Itoa(value)
}

func StarsLabel(lang Lang, value int) string {
	return strconv.Itoa(value) + " " + StarsWord(lang, value)
}

func OptionalHandle(username string) string {
	if username == "" {
		return ""
	}
	return "@" + username
}
