package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dangdennis/crossing/common/db"
)

// HelpCommand handles !help
func HelpCommand(client *db.PrismaClient, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID,
		`
Available commands:
!raid - check out the today's event.
!join - join the active raid.
!action - confirm that you've completed your daily task, and take part in the struggle!
!help - get a list of all available commands.
		`,
	)
}

// BombCommand handles !bomb
func BombCommand(client *db.PrismaClient, s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !bomb")
	// query for user rows
	// soft_delete user rows
	// soft delete avatar rows
	// soft delete avatarOnRaids rows
}
