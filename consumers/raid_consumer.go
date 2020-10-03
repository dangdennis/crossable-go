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

	raidMember, err := getAvatarRaidMembership(db, avatar, raid)
	if err != nil {
		log.Error("user is not a part of the raid", zap.Error(err), zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))
		return
	}

	currentEvent, err := getCurrentEventInStory(db, raid.Story())
	if err != nil {
		log.Error("failed to find current event for story", zap.Error(err), zap.Int("storyID", currentEvent.ID))
		return
	}

	actionMessage, err := getActionMessageForEventAndRaidMember(db, currentEvent, raidMember)
	if err != nil {
		log.Error("failed to find a message for the action",
			zap.Error(err),
			zap.Int("eventID", currentEvent.ID),
			zap.String("messageType", messages.MessageTypeActionSingle.String()),
			zap.Int("sequence", raidMember.Position),
		)
		return
	}

	err = createAvatarEventAction(db, currentEvent, avatar)
	if err != nil {
		log.Error("failed to create avatar event", zap.Error(err))
		return
	}

	// send the action's message to the user
	_, err = s.ChannelMessageSend(m.ChannelID, actionMessage.Content)
	if err != nil {
		log.Error("failed to send message", zap.Error(err))
		return
	}
}

func getAvatarRaidMembership(db *prisma.PrismaClient, avatar prisma.AvatarModel, raid prisma.RaidModel) (prisma.AvatarsOnRaidsModel, error) {
	// verify that the avatar is a part of the raid
	avatarOnRaids, err := db.AvatarsOnRaids.FindMany(
		prisma.AvatarsOnRaids.AvatarID.Equals(avatar.ID),
		prisma.AvatarsOnRaids.RaidID.Equals(raid.ID),
	).Take(1).Exec(context.Background())
	if err != nil {
		return prisma.AvatarsOnRaidsModel{}, fmt.Errorf("failed to avatar raid membership. err=%w", err)
	}

	if len(avatarOnRaids) == 0 {
		return prisma.AvatarsOnRaidsModel{}, fmt.Errorf("user is not a part of the raid")
	}

	return avatarOnRaids[0], nil
}

// createAvatarEventAction logs the avatar's action for the event
func createAvatarEventAction(db *prisma.PrismaClient, currentEvent prisma.EventModel, avatar prisma.AvatarModel) error {
	_, err := db.Action.CreateOne(
		prisma.Action.Event.Link(
			prisma.Event.ID.Equals(currentEvent.ID),
		),
		prisma.Action.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar.ID)),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create action for avatar. err=%w", err)
	}

	return nil
}

// getActionMessageForEventAndRaidMember find the relevant player message for the avatar.
// If an action message does not exist for the current raid member, use the last one in the message sequence as the default.
func getActionMessageForEventAndRaidMember(db *prisma.PrismaClient, currentEvent prisma.EventModel, raidMember prisma.AvatarsOnRaidsModel) (prisma.MessageModel, error) {
	manyMessages, err := db.Message.FindMany(
		prisma.Message.EventID.Equals(currentEvent.ID),
		prisma.Message.Type.Equals(messages.MessageTypeActionSingle.String()),
		prisma.Message.Sequence.Equals(raidMember.Position),
		prisma.Message.Default.Equals(false),
	).Take(1).Exec(context.Background())
	if err != nil {
		defaultMessages, err := db.Message.FindMany(
			prisma.Message.EventID.Equals(currentEvent.ID),
			prisma.Message.Type.Equals(messages.MessageTypeActionSingle.String()),
			prisma.Message.Default.Equals(true),
		).Take(1).Exec(context.Background())
		if err != nil {
			return prisma.MessageModel{}, fmt.Errorf("failed to find a message for the action. err=%w", err)
		}

		if len(defaultMessages) == 0 {
			return prisma.MessageModel{}, fmt.Errorf("failed to find a default message for avatar action")
		}

		return defaultMessages[0], nil
	}

	if len(manyMessages) == 0 {
		return prisma.MessageModel{}, fmt.Errorf("failed to find a message for avatar action")
	}

	return manyMessages[0], nil
}

// getCurrentEventInStory gets the next event in the story sequence that hasn't occurred yet
func getCurrentEventInStory(db *prisma.PrismaClient, story prisma.StoryModel) (prisma.EventModel, error) {
	events, err := db.Event.FindMany(
		prisma.Event.StoryID.Equals(story.ID),
		prisma.Event.Occurred.Equals(false),
	).OrderBy(
		prisma.Event.Sequence.Order(prisma.ASC),
	).Take(1).Exec(context.Background())

	if err != nil {
		return prisma.EventModel{}, fmt.Errorf("failed to find events for story. err=%w", err)
	}

	if len(events) == 0 {
		return prisma.EventModel{}, fmt.Errorf("no active event found for current storyline. err=%w", err)
	}

	return events[0], nil
}
