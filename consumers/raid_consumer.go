package consumers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/libs/logger"
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
	avatarOnRaid, err := db.AvatarsOnRaids.FindMany(
		prisma.AvatarsOnRaids.AvatarID.Equals(avatar.ID),
		prisma.AvatarsOnRaids.RaidID.Equals(raid.ID),
	).Exec(context.Background())

	if len(avatarOnRaid) == 0 {
		log.Error("user is not a part of the raid", zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))
		return
	}

	// get the next event in the story sequence that hasn't occurred yet

	// find the relevant action for the user

	// complete the action

	// send the action's message to the user
}
