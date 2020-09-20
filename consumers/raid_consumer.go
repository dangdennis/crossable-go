package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// RaidCommand handles !raid
func RaidCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !raid")
}

// JoinCommand handles !raid
func JoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !join")
}

// AttackCommand handles !raid
func AttackCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !attack")
}
