package consumers

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/dg"
	"github.com/dangdennis/crossing/common/logger"
	"github.com/dangdennis/crossing/common/repositories/messages"
	"github.com/dangdennis/crossing/common/repositories/raids"
	"github.com/dangdennis/crossing/common/repositories/stories"
	"github.com/dangdennis/crossing/common/repositories/users"
	"github.com/dangdennis/crossing/common/services/auth"
)

// RaidCommand handles !raid
func RaidCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := prisma.Client()
	log := logger.GetLogger()

	raid, err := raids.FindLatestActiveRaid(db)
	if err != nil {
		log.Error("failed to find an active raid", zap.Error(err))
		dg.ChannelMessageSend(s, m.ChannelID, "No active raid this week.")
		return
	}

	currentEvent, err := stories.GetCurrentEventInStory(db, raid.Story())
	if err != nil {
		log.Error("failed to find current event for story", zap.Error(err), zap.Int("storyID", currentEvent.ID), zap.Int("raidID", raid.ID))
		return
	}

	var message string

	eventName, ok := currentEvent.Name()
	if ok {
		message = message + eventName
	}

	introMessage, err := stories.GetEventIntroMessage(db, currentEvent)
	if err != nil {
		log.Error("failed to get intro message for event story", zap.Error(err), zap.Int("eventID", currentEvent.ID), zap.Int("raidID", raid.ID))
		return
	}

	message = fmt.Sprintf("%s\n\n%s", message, introMessage.Content)

	if currentEvent.Sequence == 1 || currentEvent.Sequence == 2 {
		message = fmt.Sprintf("%s\n\n%s", message, "`!join` to join mission!\n\n`!action` to take action.")
	}

	dg.ChannelMessageSend(s, m.ChannelID, message)
}

// JoinCommand handles !join
func JoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	raidMember, err := raids.GetAvatarRaidMembership(db, avatar, raid)
	if err == nil && raidMember.AvatarID > 0 {
		dg.ChannelMessageSend(s, m.ChannelID, "You've already joined the raid.")
		return
	}

	_, err = raids.JoinRaid(db, raid, avatar)
	if err != nil {
		log.Error("failed to add avatar to raid", zap.Error(err))
		return
	}

	username := m.Author.Username
	if username == "" {
		dg.ChannelMessageSend(s, m.ChannelID, fmt.Sprintf("A new member has joined the raid!"))
		return
	}

	dg.ChannelMessageSend(s, m.ChannelID, fmt.Sprintf("%s has joined the raid!", username))
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

	raidMember, err := raids.GetAvatarRaidMembership(db, avatar, raid)
	if err != nil {
		log.Error("user is not a part of the raid", zap.Error(err), zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))
		dg.ChannelMessageSend(s, m.ChannelID, "You're not part of the raid yet. You can `!join` within the first two days.")
		return
	}

	currentEvent, err := stories.GetCurrentEventInStory(db, raid.Story())
	if err != nil {
		log.Error("failed to find current event for story", zap.Error(err), zap.Int("storyID", currentEvent.ID))
		return
	}

	actionMessage, err := stories.GetActionMessageForEventAndRaidMember(db, currentEvent, raidMember)
	if err != nil {
		log.Error("failed to find a message for the action",
			zap.Error(err),
			zap.Int("eventID", currentEvent.ID),
			zap.String("messageType", messages.MessageTypeActionSingle.String()),
			zap.Int("sequence", raidMember.Position),
		)
		return
	}

	err = stories.CreateAvatarEventAction(db, currentEvent, avatar)
	if err != nil {
		log.Error("failed to create avatar action", zap.Error(err))
		dg.ChannelMessageSend(s, m.ChannelID, "You've already performed your action today.")
		return
	}

	err = users.AwardTokens(db, user.ID, 2)
	if err != nil {
		log.Error("failed to award user tokens", zap.Error(err), zap.Int("userID", user.ID))
	}

	// send the action's message to the user
	dg.ChannelMessageSend(s, m.ChannelID, fmt.Sprintf("%s\n\n%s", actionMessage.Content, "+2 tokens for you!"))
}

// IntroCommand handles admin-only !intro command
func IntroCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !auth.IsAdmin(m.Author.ID) {
		return
	}

	RaidCommand(s, m)
}

// OutroCommand handles admin-only !outro command
func OutroCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !auth.IsAdmin(m.Author.ID) {
		return
	}

	db := prisma.Client()
	log := logger.GetLogger()

	raid, err := raids.FindLatestActiveRaid(db)
	if err != nil {
		log.Error("failed to find an active raid", zap.Error(err))
		dg.ChannelMessageSend(s, m.ChannelID, "No active raid this week.")
		return
	}

	currentEvent, err := stories.GetCurrentEventInStory(db, raid.Story())
	if err != nil {
		log.Error("failed to find current event for story", zap.Error(err), zap.Int("storyID", currentEvent.ID), zap.Int("raidID", raid.ID))
		return
	}

	outroMessage, err := stories.GetEventOutroMessage(db, currentEvent)
	if err != nil {
		log.Error("failed to get outro message for event story", zap.Error(err), zap.Int("eventID", currentEvent.ID), zap.Int("raidID", raid.ID))
		return
	}

	var engagedUsersMsg string
	actions, err := db.Action.FindMany(
		prisma.Action.EventID.Equals(currentEvent.ID),
		).With(
			prisma.Action.Avatar.Fetch().With(
				prisma.Avatar.User.Fetch(),),
			).Exec(context.Background())
	if err != nil {
		log.Error("failed to find the actions performed for event", zap.Error(err), zap.Int("eventID", currentEvent.ID))
	} else {
		engagedUsersMsg = ""
		for _, action := range actions {
			username, ok := action.Avatar().User().DiscordUsername()
			if ok {
				engagedUsersMsg = engagedUsersMsg + username + ", "
			}
		}
		engagedUsersMsg = strings.TrimRight(strings.Trim(engagedUsersMsg, " "), ",")
		engagedUsersMsg = engagedUsersMsg + ": You're our heroes."
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n\n%s",outroMessage.Content, engagedUsersMsg))
	if err != nil {
		fmt.Println(err)
		return
	}
}
