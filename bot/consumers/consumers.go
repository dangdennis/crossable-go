package consumers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/users"
)

// HandleMessageCreate consumes Discord MessageCreate events
func HandleMessageCreate(client *db.PrismaClient, s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	err := initUser(client, m)
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
		RaidCommand(client, s, m)
	case "!join":
		JoinCommand(client, s, m)
	case "!action":
		ActionCommand(client, s, m)
	case "!help":
		HelpCommand(client, s, m)
	case "!bomb":
		BombCommand(client, s, m)
	case "!intro":
		IntroCommand(client, s, m)
	case "!outro":
		OutroCommand(client, s, m)
	}
}

// initUser creates a new user, avatar, and wallet
func initUser(client *db.PrismaClient, m *discordgo.MessageCreate) error {
	// Consider hardening this with an additional cache layer. Check the LRU cache for a discord user id that's recently messaged the channel
	if !strings.HasPrefix(m.Content, "!") {
		return nil
	}

	_, err := users.FindUserByDiscordID(client, m.Author.ID)
	if err == nil {
		fmt.Println("user already exists")
		return nil
	}

	fmt.Println("initializing new user")

	user, err := users.CreateUser(client, users.UserAttrs{DiscordUserID: m.Author.ID, DiscordUsername: &m.Author.Username})
	if err != nil {
		return err
	}

	_, err = users.CreateAvatar(client, user.ID)
	if err != nil {
		return err
	}

	_, err = users.CreateWallet(client, user.ID)
	if err != nil {
		return err
	}

	return nil
}
