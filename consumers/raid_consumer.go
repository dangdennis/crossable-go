package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// RaidCommand handles !raid
func RaidCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !raid")
}

// JoinCommand handles !join
func JoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !join")
}

// ActionCommand handles !action
func ActionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !action")
}
