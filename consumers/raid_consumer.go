package consumers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/repositories"
)

// RaidCommand handles !raid
func RaidCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	raid, err := repositories.FindWeeklyActiveRaid(db.Client())
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
	// query active raid
	// query for user
	// query for their avatar
	// add their avatar to avatarOnRaids for active raid
}

// ActionCommand handles !action
func ActionCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("handling !action")
	// query active raid
	// query for user
	// query for avatar
	// perform action
	// - reveal next story segment
	// - increase raid's completion progress
}
