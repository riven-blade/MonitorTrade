package config

import (
	"os"
	"strconv"
)

type Config struct {
	Redis             RedisConfig `json:"redis"`          // Redis configuration
	TelegramToken     string      `json:"telegram_token"` // Telegram configuration
	TelegramId        int64       `json:"telegram_id"`    // Telegram configuration
	FundingRate       float64     `json:"funding_rate"`   // Funding rate threshold
	BotBaseUrl        string      `json:"freqtrade_base_url"`
	BotUsername       string      `json:"bot_username"`
	BotPasswd         string      `json:"bot_passwd"`
	BotAdjustEntryTag string      `json:"bot_adjust_entry_tag"`
}

type RedisConfig struct {
	Addr      string `json:"addr"`
	Password  string `json:"password"` // Redis password
	DB        int    `json:"db"`       // Redis database number
	KeyExpire int    `json:"key_expire"`
}

// getEnvString 获取字符串环境变量，如不存在则使用默认值
func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt 获取整数类型环境变量，如不存在或格式错误则使用默认值
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvFloat64
func getEnvFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatVal
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	config := &Config{
		Redis: RedisConfig{
			Addr:      getEnvString("REDIS_ADDR", "redis:6379"),
			Password:  getEnvString("REDIS_PASSWORD", ""),
			DB:        getEnvInt("REDIS_DB", 0),
			KeyExpire: getEnvInt("KEY_EXPIRE", 2592000),
		},
		TelegramToken:     getEnvString("TELEGRAM_TOKEN", ""),
		TelegramId:        int64(getEnvInt("TELEGRAM_ID", 0)),
		FundingRate:       getEnvFloat64("FUNDING_RATE", -0.1),
		BotBaseUrl:        getEnvString("BOT_BASE_URL", "http://127.0.0.1:8080"),
		BotUsername:       getEnvString("BOT_USER_NAME", ""),
		BotPasswd:         getEnvString("BOT_PASSWD", ""),
		BotAdjustEntryTag: getEnvString("BOT_ADJUST_ENTRY_TAG", "grind_3_entry"),
	}
	return config
}
