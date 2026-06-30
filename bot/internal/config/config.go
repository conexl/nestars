package config

import "os"

type Links struct {
	Support     string
	NewsChannel string
}

type Config struct {
	Telegram struct {
		Token string
	}
	BackendURL      string
	MockBackend     bool
	MongoURI        string
	MongoDB         string
	Links           Links
	MiniAppURL        string
	MiniAppPort       string
	MiniAppStaticPath string
	DisableMiniAppAuth bool
}

// Cfg глобальная ссылка на актуальную конфигурацию.
// Используется UI-клавиатурами, которым неудобно пробрасывать Config отдельно.
var Cfg *Config

func NewConfig() *Config {
	c := &Config{}
	c.Telegram.Token = os.Getenv("BOT_TOKEN")
	c.BackendURL = os.Getenv("BACKEND_URL")
	c.MockBackend = os.Getenv("MOCK_BACKEND") == "true"
	c.MongoURI = getEnv("MONGO_URI", "mongodb://localhost:27017")
	c.MongoDB = getEnv("MONGO_DB", "nestars_bot")
	c.Links.Support = getEnv("SUPPORT_URL", "")
	c.Links.NewsChannel = getEnv("NEWS_CHANNEL_URL", "")
	c.MiniAppURL = getEnv("MINIAPP_URL", "")
	c.MiniAppPort = getEnv("MINIAPP_PORT", "8080")
	c.MiniAppStaticPath = getEnv("MINIAPP_STATIC_PATH", "./miniapp/dist")
	c.DisableMiniAppAuth = getEnv("MINIAPP_DISABLE_AUTH", "") == "true"
	Cfg = c
	return c
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
