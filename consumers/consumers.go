package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/dangdennis/crossing/db"
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

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func initUser(m *discordgo.MessageCreate) error {
	_, err := repositories.FindUserByDiscordID(db.Client(), m.Author.ID)
	if err == nil {
		fmt.Println("user already exists")
		return nil
	} else {
		fmt.Println("initializing new user")
	}

	user, err := repositories.CreateUser(db.Client(), repositories.UserAttrs{DiscordUserID: m.Author.ID})
	if err != nil {
		return err
	}

	_, err = repositories.CreateAvatar(db.Client(), user.ID)
	if err != nil {
		return err
	}

	return nil
}
