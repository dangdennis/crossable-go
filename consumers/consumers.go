package consumers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/repositories"
)

// MessageCreate consumes Discord MessageCreate events
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// fmt.Println(s)
	fmt.Println(m.Author.ID)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	err := initUser(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch m.Content {
	case "!ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "!pong":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "!raid":
		RaidCommand(s, m)
	case "!join":
		JoinCommand(s, m)
	case "!action":
		ActionCommand(s, m)
	case "!help":
		HelpCommand(s, m)
	case "!bomb":
		BombCommand(s, m)
	}
}

// initUser creates a new user account if the current account is not found
func initUser(m *discordgo.MessageCreate) error {
	// Consider hardening this with an additional cache layer. Check the LRU cache for a discord user id that's recently messaged the channel
	if !strings.HasPrefix(m.Content, "!") {
		return nil
	}

	_, err := repositories.FindUserByDiscordID(prisma.Client(), m.Author.ID)
	if err == nil {
		fmt.Println("user already exists")
		return nil
	} else {
		fmt.Println("initializing new user")
	}

	user, err := repositories.CreateUser(prisma.Client(), repositories.UserAttrs{DiscordUserID: m.Author.ID})
	if err != nil {
		return err
	}

	_, err = repositories.CreateAvatar(prisma.Client(), user.ID)
	if err != nil {
		return err
	}

	return nil
}
