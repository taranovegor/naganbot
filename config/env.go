package config

import "os"

const (
	TelegramBotToken = "TELEGRAM_BOT_TOKEN"
	DatabaseDsn      = "DATABASE_DSN"
	DrumCapacity     = "DRUM_CAPACITY"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}
