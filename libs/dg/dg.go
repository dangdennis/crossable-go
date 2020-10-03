package dg

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/dangdennis/crossing/libs/logger"
)

// ChannelMessageSend sends a message to Discord and logs any error
func ChannelMessageSend(s *discordgo.Session, channelID string, message string) {
	log := logger.GetLogger()
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Error("failed to send message", zap.Error(err))
	}
}
