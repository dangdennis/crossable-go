package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HelpCommand handles !help
func HelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID,
		`
Available commands:
!raid - status on the weekly raid.
!join - join the week's raid.
!action - confirm that you've completed your daily task, and take part in the struggle!
!help - get a list of all available commands.

Real serious commands:
!bomb - deletes all your data.
		`,
	)
}

// BombCommand handles !bomb
func BombCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !bomb")
}
