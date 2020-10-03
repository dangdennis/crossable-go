package consumers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/libs/logger"
	"github.com/dangdennis/crossing/repositories/messages"
	"github.com/dangdennis/crossing/repositories/raids"
	"github.com/dangdennis/crossing/repositories/users"
)

// RaidCommand handles !raid
func RaidCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	raid, err := raids.FindLatestActiveRaid(prisma.Client())
	if err != nil {
		fmt.Println(err)
		_, err := s.ChannelMessageSend(m.ChannelID, `No active raid this week.`)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	raidBossesOnRaids := raid.RaidBossesOnRaids()
	if len(raidBossesOnRaids) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Raid bosses are in hiding still."))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	boss := raidBossesOnRaids[0].RaidBoss()

	if raid.CompletionProgress < 1 {
		health := (1 - raid.CompletionProgress) * 100
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %.2f%% HP left!", boss.Name, health))
		if err != nil {
			fmt.Println(err)
			return
		}
		return

	}

	successMsg := fmt.Sprintf("This week's raid boss, %s, has been defeated. \nJoin next week's raid!", boss.Name)
	_, err = s.ChannelMessageSend(m.ChannelID, successMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// JoinCommand handles !join
func JoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !join")
	db := prisma.Client()
	log := logger.GetLogger()

	raid, err := raids.FindLatestActiveRaid(db)
	if err != nil {
		log.Error("failed to get weekly active raid", zap.Error(err))
		return
	}

	user, err := users.FindUserByDiscordID(db, m.Author.ID)
	if err != nil {
		log.Error("failed to get weekly active raid", zap.Error(err))
		return
	}

	avatar, ok := user.Avatar()
	if !ok {
		log.Error("user does not have avatar", zap.Int("userID", user.ID), zap.Error(err))
		return
	}

	_, err = raids.JoinRaid(db, raid, avatar)
	if err != nil {
		log.Error("failed to add avatar to raid", zap.Error(err))
		return
	}

	username := m.Author.Username
	if username == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A new member has joined the raid!"))
		if err != nil {
			log.Error("failed to send message", zap.Error(err))
		}
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has joined the raid!", username))
	if err != nil {
		log.Error("failed to send message", zap.Error(err))
	}
}

// ActionCommand handles !action
func ActionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !action")
	db := prisma.Client()
	log := logger.GetLogger()

	raid, err := raids.FindLatestActiveRaid(db)
	if err != nil {
		log.Error("failed to get weekly active raid", zap.Error(err))
		return
	}

	user, err := users.FindUserByDiscordID(db, m.Author.ID)
	if err != nil {
		log.Error("failed to get weekly active raid", zap.Error(err))
		return
	}

	avatar, ok := user.Avatar()
	if !ok {
		log.Error("user does not have avatar", zap.Int("userID", user.ID), zap.Error(err))
		return
	}

	// verify that the avatar is a part of the raid
	avatarOnRaids, err := db.AvatarsOnRaids.FindMany(
		prisma.AvatarsOnRaids.AvatarID.Equals(avatar.ID),
		prisma.AvatarsOnRaids.RaidID.Equals(raid.ID),
	).Take(1).Exec(context.Background())

	if len(avatarOnRaids) == 0 {
		log.Error("user is not a part of the raid", zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))
		return
	}

	raidMember := avatarOnRaids[0]
	currentStory := raid.Story()

	// get the next event in the story sequence that hasn't occurred yet
	events, err := db.Event.FindMany(
		prisma.Event.StoryID.Equals(currentStory.ID),
		prisma.Event.Occurred.Equals(false),
	).OrderBy(
		prisma.Event.Sequence.Order(prisma.ASC),
	).Take(1).Exec(context.Background())
	if err != nil {
		log.Error("failed to find events for story", zap.Error(err), zap.Int("storyID", currentStory.ID))
		return
	}

	if len(events) == 0 {
		log.Error("no active event found for current storyline", zap.Int("storyID", currentStory.ID))
		return
	}

	currentEvent := events[0]

	// find the relevant player message for the avatar
	msgs, err := db.Message.FindMany(
		prisma.Message.EventID.Equals(currentEvent.ID),
		prisma.Message.Type.Equals(messages.MessageTypeActionSingle.String()),
		prisma.Message.Sequence.Equals(raidMember.Position),
	).Take(1).Exec(context.Background())
	if err != nil {
		log.Error("failed to find a message for the action",
			zap.Error(err),
			zap.Int("eventID", currentEvent.ID),
			zap.String("messageType", messages.MessageTypeActionSingle.String()),
			zap.Int("sequence", raidMember.Position),
		)
	}
	if len(msgs) == 0 {
		log.Error("failed to find a message for avatar action")
		return
	}

	actionMessage := msgs[0]

	// log the avatar's action for the event
	_, err = db.Action.CreateOne(
		prisma.Action.Event.Link(
			prisma.Event.ID.Equals(currentEvent.ID),
		),
		prisma.Action.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar.ID)),
	).Exec(context.Background())
	if err != nil {
		log.Error("failed to create action for avatar", zap.Error(err), zap.Int("avatarID", avatar.ID))
		return
	}

	// send the action's message to the user
	_, err = s.ChannelMessageSend(m.ChannelID, actionMessage.Content)
	if err != nil {
		log.Error("failed to send message", zap.Error(err))
	}
}
