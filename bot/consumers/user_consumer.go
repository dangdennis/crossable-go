package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HelpCommand handles !help
func HelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !help")
}

// BombCommand handles !bomb
func BombCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !bomb")
}
