package raids

import (
	"context"
	"fmt"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/libs/logger"
	"go.uber.org/zap"
)

// FindWeeklyActiveRaid gets the active raid of the week and its raid bosses
// TODO Pass time boundaries (start and end of a week) for better testability
func FindWeeklyActiveRaid(db *prisma.PrismaClient) (r prisma.RaidModel, err error) {
	raids, err := db.Raid.FindMany(
		prisma.Raid.Active.Equals(true),
	).With(
		prisma.Raid.RaidBossesOnRaids.Fetch().Take(5).With(
			prisma.RaidBossesOnRaids.RaidBoss.Fetch()),
	).OrderBy(
		prisma.Raid.StartTime.Order(prisma.DESC),
	).Take(1).Exec(context.Background())
	if err != nil {
		return r, err
	}

	if len(raids) == 0 {
		return r, fmt.Errorf("failed to find the most recently started active raid")
	}

	return raids[0], nil
}

// JoinRaid adds a user's avatar to a raid as a member
func JoinRaid(db *prisma.PrismaClient, raid prisma.RaidModel, avatar prisma.AvatarModel) error {
	_, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar.ID),
		),
	).Exec(context.Background())
	if err != nil {
		return err
	}

	logger.GetLogger().Info("successfully added avatar to raid", zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))

	return nil
}
