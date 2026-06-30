package i18n

import "fmt"

func MainMenu(lang Lang) string {
	if lang == RU {
		return `<tg-emoji emoji-id="5994750571041525522">👋</tg-emoji> <b>Привет!</b>

<blockquote><tg-emoji emoji-id="6032937473162614352">🎁</tg-emoji> Здесь ты можешь приобрести Telegram звезды и Premium подписки без 
верификации KYC и дешевле чем в приложении.
</blockquote>
`
	}

	return `<tg-emoji emoji-id="5994750571041525522">👋</tg-emoji> <b>Hello!</b>

<blockquote><tg-emoji emoji-id="6032937473162614352">🎁</tg-emoji> Here you can buy Telegram stars and Premium subscriptions without 
KYC verification and cheaper than in the app.
</blockquote>
`
}

func StarsWho(lang Lang) string {
	if lang == RU {
		return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Выберите получателя</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Пожалуйста, выберите кому хотите отправить звезды. 
Можно подарить другу введя его @username.
</blockquote>
`
	}

	return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Select recipient</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Choose who you want to send stars to. 
You can gift them to a friend by entering their @username.
</blockquote>
`
}

func PremiumWho(lang Lang) string {
	if lang == RU {
		return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Выберите получателя</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Пожалуйста, выберите кому хотите подарить Premium.
Можно указать друга через @username.
</blockquote>
`
	}

	return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Select recipient</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Choose who you want to gift Premium to.
You can enter a friend's @username.
</blockquote>
`
}

func WithdrawWho(lang Lang) string {
	if lang == RU {
		return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Выберите получателя</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Укажите кому хотите вывести звезды.
Можно ввести @username получателя.
</blockquote>
`
	}

	return `<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Select recipient</b>

<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Choose who you want to withdraw stars to.
You can enter the recipient's @username.
</blockquote>
`
}

func StarsFlow(lang Lang, recipient, product string) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Оформление заказа</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Получатель:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Количество звезд:</b> %s</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Вы можете купить любое количество звёзд.
Нажмите «Ввести количество» и укажите число — любое.
</blockquote>
`, recipient, product)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Create order</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Recipient:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Stars amount:</b> %s</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> You can buy any number of stars.
Tap «Enter amount» and specify any number you want.
</blockquote>
`, recipient, product)
}

func PremiumFlow(lang Lang, recipient, product string) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Оформление заказа</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Получатель:</b> %s
<tg-emoji emoji-id="5956561749070057536">⌛</tg-emoji> <b>Период подписки:</b> %s</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Premium подписка дает доступ к расширенным функциям Telegram.</blockquote>
`, recipient, product)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Create order</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Recipient:</b> %s
<tg-emoji emoji-id="5956561749070057536">⌛</tg-emoji> <b>Subscription period:</b> %s</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Premium gives access to advanced Telegram features.</blockquote>
`, recipient, product)
}

func PaymentScreen(lang Lang, isPremium bool, username, product, amount string) string {
	if isPremium {
		if lang == RU {
			return fmt.Sprintf(`<tg-emoji emoji-id="5415696465945173620">🦄</tg-emoji> <b>Оплата заказа</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Пользователь:</b> %s
<tg-emoji emoji-id="5956561749070057536">⌛</tg-emoji> <b>Период подписки:</b> %s
<tg-emoji emoji-id="5778318458802409852">💰</tg-emoji> <b>Сумма:</b> %s</blockquote>

<blockquote><tg-emoji emoji-id="5927169041595634481">💳</tg-emoji> <b>Выберите способ оплаты:</b></blockquote>
`, username, product, amount)
		}

		return fmt.Sprintf(`<tg-emoji emoji-id="5415696465945173620">🦄</tg-emoji> <b>Order payment</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>User:</b> %s
<tg-emoji emoji-id="5956561749070057536">⌛</tg-emoji> <b>Subscription period:</b> %s
<tg-emoji emoji-id="5778318458802409852">💰</tg-emoji> <b>Total:</b> %s</blockquote>

<blockquote><tg-emoji emoji-id="5927169041595634481">💳</tg-emoji> <b>Choose payment method:</b></blockquote>
`, username, product, amount)
	}

	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Оплата заказа</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Пользователь:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Количество:</b> %s
<tg-emoji emoji-id="5778318458802409852">💰</tg-emoji> <b>Сумма:</b> %s</blockquote>

<blockquote><tg-emoji emoji-id="5927169041595634481">💳</tg-emoji> <b>Выберите способ оплаты:</b></blockquote>
`, username, product, amount)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="6032937473162614352">⭐</tg-emoji> <b>Order payment</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>User:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Amount:</b> %s
<tg-emoji emoji-id="5778318458802409852">💰</tg-emoji> <b>Total:</b> %s</blockquote>

<blockquote><tg-emoji emoji-id="5927169041595634481">💳</tg-emoji> <b>Choose payment method:</b></blockquote>
`, username, product, amount)
}

func PaymentReady(lang Lang, amount, orderNumber string) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="5778605968208170641">🕒</tg-emoji> <b>Счёт действителен 10 минут</b>

<blockquote><tg-emoji emoji-id="5951584964305755220">🧾</tg-emoji> Номер заказа: <b>%s</b>
<tg-emoji emoji-id="5987880246865565644">💰</tg-emoji> Сумма: <b>%s</b>

<tg-emoji emoji-id="5967591100532134862">📞</tg-emoji> В случае возникновения проблем, обращайтесь в поддержку.
</blockquote>
`, orderNumber, amount)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="5778605968208170641">🕒</tg-emoji> <b>Invoice is valid for 10 minutes</b>

<blockquote><tg-emoji emoji-id="5951584964305755220">🧾</tg-emoji> Order number: <b>%s</b>
<tg-emoji emoji-id="5987880246865565644">💰</tg-emoji> Amount: <b>%s</b>

<tg-emoji emoji-id="5967591100532134862">📞</tg-emoji> If you have any problems, please contact support.
</blockquote>
`, orderNumber, amount)
}

func PaymentSuccessBalance(lang Lang, amount, orderNumber string) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="6051321050235283584">✅</tg-emoji> <b>Заказ оплачен с баланса</b>

<blockquote><tg-emoji emoji-id="5951584964305755220">🧾</tg-emoji> Номер заказа: <b>%s</b>
<tg-emoji emoji-id="5987880246865565644">💰</tg-emoji> Сумма: <b>%s</b>

<tg-emoji emoji-id="6051321050235283584">✅</tg-emoji> Спасибо за покупку!
</blockquote>
`, orderNumber, amount)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="6051321050235283584">✅</tg-emoji> <b>Order paid from balance</b>

<blockquote><tg-emoji emoji-id="5951584964305755220">🧾</tg-emoji> Order number: <b>%s</b>
<tg-emoji emoji-id="5987880246865565644">💰</tg-emoji> Amount: <b>%s</b>

<tg-emoji emoji-id="6051321050235283584">✅</tg-emoji> Thank you for your purchase!
</blockquote>
`, orderNumber, amount)
}

func ProfileText(lang Lang, telegramID int64, createdAt string, balanceUSD, balanceTON float64, boughtStars int) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Профиль</b>
	
<blockquote><tg-emoji emoji-id="6028435952299413210">🆔</tg-emoji> ID: %d
<tg-emoji emoji-id="6142986161762607580">⚒️</tg-emoji> Регистрация: <b>%s</b>

<tg-emoji emoji-id="5769403330761593044">🏦</tg-emoji><b> Баланс: %.2f$ (%.2f TON)</b>

<tg-emoji emoji-id="5931472654660800739">📊</tg-emoji> <b>Статистика:</b>
<tg-emoji emoji-id="5318770787925113164">🟢</tg-emoji> Куплено звезд: <b>%d</b>
</blockquote>
`, telegramID, createdAt, balanceUSD, balanceTON, boughtStars)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Profile</b>
	
<blockquote><tg-emoji emoji-id="6028435952299413210">🆔</tg-emoji> ID: %d
<tg-emoji emoji-id="6142986161762607580">⚒️</tg-emoji> Registered: <b>%s</b>

<tg-emoji emoji-id="5769403330761593044">🏦</tg-emoji><b> Balance: %.2f$ (%.2f TON)</b>

<tg-emoji emoji-id="5931472654660800739">📊</tg-emoji> <b>Stats:</b>
<tg-emoji emoji-id="5318770787925113164">🟢</tg-emoji> Bought stars: <b>%d</b>
</blockquote>
`, telegramID, createdAt, balanceUSD, balanceTON, boughtStars)
}

func ReferralProgramText(lang Lang, link string, earned, invited, balance int) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="5960672896060756972">⭐</tg-emoji> <b>Реферальная программа</b>
		
<tg-emoji emoji-id="5877465816030515018">🔗</tg-emoji> <b>Ссылка:</b>
%s

<tg-emoji emoji-id="5875291072225087249">📊</tg-emoji> <b>Статистика:</b>
<blockquote><tg-emoji emoji-id="5958376256788502078">⭐</tg-emoji> <b>Заработано: %d</b>
<tg-emoji emoji-id="5877530150345641603">👥</tg-emoji> <b>Приглашено: %d</b>
<tg-emoji emoji-id="5769403330761593044">👛</tg-emoji> <b>Баланс: <tg-emoji emoji-id="5958376256788502078">⭐</tg-emoji> %d</b>
</blockquote>`, link, earned, invited, balance)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="5960672896060756972">⭐</tg-emoji> <b>Referral program</b>
		
<tg-emoji emoji-id="5877465816030515018">🔗</tg-emoji> <b>Link:</b>
%s

<tg-emoji emoji-id="5875291072225087249">📊</tg-emoji> <b>Stats:</b>
<blockquote><tg-emoji emoji-id="5958376256788502078">⭐</tg-emoji> <b>Earned: %d</b>
<tg-emoji emoji-id="5877530150345641603">👥</tg-emoji> <b>Invited: %d</b>
<tg-emoji emoji-id="5769403330761593044">👛</tg-emoji> <b>Balance: <tg-emoji emoji-id="5958376256788502078">⭐</tg-emoji> %d</b>
</blockquote>`, link, earned, invited, balance)
}

func WithdrawFlow(lang Lang, recipient, amount string, balance int) string {
	if lang == RU {
		return fmt.Sprintf(`<tg-emoji emoji-id="5960672896060756972">⭐</tg-emoji> <b>Вывод звезд</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Получатель:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Количество звезд:</b> %s
<tg-emoji emoji-id="5769403330761593044">🏦</tg-emoji> <b>Баланс:</b> %d</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Выберите получателя и количество звёзд для вывода.
</blockquote>
`, recipient, amount, balance)
	}

	return fmt.Sprintf(`<tg-emoji emoji-id="5960672896060756972">⭐</tg-emoji> <b>Withdraw stars</b>

<blockquote><tg-emoji emoji-id="6145278957039131275">👤</tg-emoji> <b>Recipient:</b> %s
<tg-emoji emoji-id="5956561749070057536">🌟</tg-emoji> <b>Stars amount:</b> %s
<tg-emoji emoji-id="5769403330761593044">🏦</tg-emoji> <b>Balance:</b> %d</blockquote>
<blockquote><tg-emoji emoji-id="5873121512445187130">❓</tg-emoji> Select recipient and stars amount for withdrawal.
</blockquote>
`, recipient, amount, balance)
}

func NotSelected(lang Lang) string {
	if lang == RU {
		return "не выбран"
	}
	return "not selected"
}

func NotSelectedValue(lang Lang) string {
	if lang == RU {
		return "не выбрано"
	}
	return "not selected"
}

func NoProducts(lang Lang) string {
	if lang == RU {
		return "⚠️ Товаров нет"
	}
	return "⚠️ No products"
}

func StarsCustomAmount(lang Lang) string {
	if lang == RU {
		return "ℹ️ Введите количество ⭐ вручную:"
	}
	return "ℹ️ Enter stars amount manually:"
}

func InvalidStarsAmount(lang Lang) string {
	if lang == RU {
		return "⚠️ Введите количество от 50 до 1.000.000 ⭐"
	}
	return "⚠️ Enter an amount between 50 and 1,000,000 ⭐"
}

func UsernameRequired(lang Lang) string {
	if lang == RU {
		return "❗ Для покупки звёзд или Premium необходимо иметь @username"
	}
	return "❗ You need a public @username to buy stars or Premium"
}

func EnterFriendUsername(lang Lang) string {
	if lang == RU {
		return "👤 Введите username друга (например @username или https://t.me/username):"
	}
	return "👤 Enter friend's username (for example @username or https://t.me/username):"
}

func InvalidUsername(lang Lang) string {
	if lang == RU {
		return "⚠️ Введите корректный username (например @username или https://t.me/username):"
	}
	return "⚠️ Enter a valid username (for example @username or https://t.me/username):"
}

func WithdrawAmount(lang Lang) string {
	if lang == RU {
		return "ℹ️ Введите количество звёзд для вывода:"
	}
	return "ℹ️ Enter stars amount for withdrawal:"
}

func InvalidWithdrawStars(lang Lang) string {
	if lang == RU {
		return "⚠️ Введите корректное количество звёзд для вывода"
	}
	return "⚠️ Enter a valid stars amount for withdrawal"
}

func CancelOrderAlert(lang Lang) string {
	if lang == RU {
		return "😔 Заказ успешно отменен!"
	}
	return "😔 Order has been canceled successfully!"
}

func GenericCatalogError(lang Lang) string {
	if lang == RU {
		return "❗ Ошибка загрузки каталога"
	}
	return "❗ Failed to load catalog"
}

func NoPremiumPackages(lang Lang) string {
	if lang == RU {
		return "⚠️ В каталоге нет доступных пакетов Premium"
	}
	return "⚠️ No Premium packages available in catalog"
}

func StartRequired(lang Lang) string {
	if lang == RU {
		return "❗ Сначала нажмите /start"
	}
	return "❗ Please press /start first"
}

func UUIDError(lang Lang) string {
	if lang == RU {
		return "❗ Ошибка UUID"
	}
	return "❗ UUID error"
}

func CreateUserError(lang Lang) string {
	if lang == RU {
		return "❗ Ошибка при создании пользователя"
	}
	return "❗ Failed to create user"
}

func ProfileLoadError(lang Lang) string {
	if lang == RU {
		return "❗ Ошибка загрузки профиля"
	}
	return "❗ Failed to load profile"
}

func PaymentOrderIDMissing(lang Lang) string {
	if lang == RU {
		return "❗ Не найден OrderID для оплаты"
	}
	return "❗ OrderID not found for payment"
}

func SelectRecipientAndPackage(lang Lang) string {
	if lang == RU {
		return "Выберите получателя и пакет ⭐"
	}
	return "Select recipient and package ⭐"
}

func FillRecipientAndAmount(lang Lang) string {
	if lang == RU {
		return "Заполните получателя и количество звёзд"
	}
	return "Fill in recipient and stars amount"
}

func WithdrawSuccess(lang Lang) string {
	if lang == RU {
		return "✅ Звезды успешно выведены"
	}
	return "✅ Stars withdrawn successfully"
}
