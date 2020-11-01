package dg

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/dangdennis/crossing/common/logger"
)

// ChannelMessageSend sends a message to a channel
func ChannelMessageSend(s *discordgo.Session, channelID string, message string) {
	log := logger.GetLogger()
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Error("failed to send message", zap.Error(err))
	}
}

// DirectMessageSend sends a direct message to a user
func DirectMessageSend(s *discordgo.Session, recipientID string, message string) {
	log := logger.GetLogger()
	dmChannel, err := s.UserChannelCreate(recipientID)
	if err != nil {
		log.Error("failed to get user's direct channel", zap.Error(err))
		return
	}

	ChannelMessageSend(s, dmChannel.ID, message)
}
