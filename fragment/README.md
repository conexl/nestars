# Fragment API REST Wrapper

Stateless REST API на Python поверх библиотеки [`fragment-api-py`](https://github.com/S1qwy/fragment-api-py.git) с автоматической документацией Swagger / OpenAPI.

## Возможности

- **Full mode по умолчанию** — передавайте `x-fragment-cookies` для полного доступа ко всем методам Fragment.
- **No-Cookie / No KYC mode** — не передавайте cookies, и покупки Stars / Premium / TON Ads будут работать через hosted REST API библиотеки.
- **Все методы библиотеки** — wallet, profile, marketplace, bids, purchases, giveaways, batch operations, NFT transfer/withdrawal, anonymous numbers, assignment, auctions, raw API call.
- **Stateless** — никакие сиды или cookies не хранятся на сервере, они передаются в заголовках каждого запроса.

## Установка

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## Запуск

```bash
source .venv/bin/activate
uvicorn main:app --reload --host 0.0.0.0 --port 8000
```

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
- OpenAPI JSON: http://localhost:8000/openapi.json

## Авторизация

Каждый запрос должен содержать как минимум:

| Заголовок | Описание | Обязательный |
|-----------|----------|--------------|
| `x-fragment-seed` | Мнемоническая фраза кошелька TON (12/18/24 слов) | Да |
| `x-fragment-wallet-version` | `V4R2` или `V5R1` (по умолчанию `V5R1`) | Нет |
| `x-fragment-cookies` | JSON-объект или строка `key=val;key2=val2` с Fragment-сессией | Нет |
| `x-fragment-api-key` | Пользовательский ключ TONAPI | Нет |
| `x-fragment-timeout` | Таймаут HTTP-запроса | Нет |
| `x-fragment-stats-enabled` | Включить анонимную телеметрию (`true`/`false`) | Нет |

## Примеры

### 1. Получить информацию о кошельке

```bash
curl -X GET "http://localhost:8000/wallet" \
  -H "x-fragment-seed: word1 word2 ... word24"
```

### 2. Купить Stars

```bash
curl -X POST "http://localhost:8000/purchase/stars" \
  -H "x-fragment-seed: word1 word2 ... word24" \
  -H "Content-Type: application/json" \
  -d '{"username": "@durov", "amount": 100}'
```

### 3. Batch-покупка

```bash
curl -X POST "http://localhost:8000/purchase/batch" \
  -H "x-fragment-seed: word1 word2 ... word24" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"type": "stars", "username": "@durov", "amount": 100},
      {"type": "premium", "username": "@telegram", "months": 6}
    ],
    "payment_method": "ton"
  }'
```

### 4. Полный режим (с cookies)

```bash
curl -X GET "http://localhost:8000/profile" \
  -H "x-fragment-seed: word1 word2 ... word24" \
  -H "x-fragment-cookies: '{\"stel_ssid\":\"...\",\"stel_dt\":\"...\",\"stel_token\":\"...\",\"stel_ton_token\":\"...\"}'"
```

### 5. Аутентификация (получить cookies)

```bash
curl -X POST "http://localhost:8000/auth/authenticate" \
  -H "Content-Type: application/json" \
  -d '{"seed": "word1 word2 ... word24", "wallet_version": "V5R1"}'
```

## Структура проекта

```
.
├── config.py            # Настройки FastAPI
├── dependencies.py      # Dependency для создания FragmentClient
├── generate_openapi.py  # Генерация openapi.json / swagger.yaml
├── main.py              # FastAPI endpoints
├── openapi.json         # OpenAPI 3.1 спецификация (JSON)
├── requirements.txt
├── schemas.py           # Pydantic модели
└── swagger.yaml         # OpenAPI 3.1 спецификация (YAML)
```

## Перегенерация Swagger

```bash
source .venv/bin/activate
python generate_openapi.py
```

## Безопасность

- Не храните seed и cookies в `.env` этого сервиса — передавайте их в заголовках.
- Размещайте API за HTTPS и ограничьте доступ (например, через reverse proxy / VPN).
- Отключите телеметрию библиотеки, если нужно: `x-fragment-stats-enabled: false`.

## Лицензия

MIT
