package env

import "os"

// GetDiscordBotKey gets the bot api token
func GetDiscordBotKey() string {
	return os.Getenv("DISCORD_BOT_KEY")
}
