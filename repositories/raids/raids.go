package raids

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/libs/logger"
	"github.com/dangdennis/crossing/repositories/story"
)

// CreateRaid creates a new raid
func CreateRaid(db *prisma.PrismaClient) (r prisma.RaidModel, err error) {
	newStory, err := story.CreateStory(db)
	if err != nil {
		return r, err
	}

	return db.Raid.CreateOne(
		prisma.Raid.Story.Link(
			prisma.Story.ID.Equals(newStory.ID),
		),
	).Exec(context.Background())
}

// FindLatestActiveRaid gets the active raid of the week and its raid bosses
// TODO Pass time boundaries (start and end of a week) for better testability
func FindLatestActiveRaid(db *prisma.PrismaClient) (r prisma.RaidModel, err error) {
	raids, err := db.Raid.FindMany(
		prisma.Raid.Active.Equals(true),
	).With(
		prisma.Raid.RaidBossesOnRaids.Fetch().Take(5).With(
			prisma.RaidBossesOnRaids.RaidBoss.Fetch()),
		prisma.Raid.AvatarsOnRaids.Fetch().With(
			prisma.AvatarsOnRaids.Avatar.Fetch(),
		),
		prisma.Raid.Story.Fetch(),
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
func JoinRaid(db *prisma.PrismaClient, raid prisma.RaidModel, avatar prisma.AvatarModel) (prisma.AvatarsOnRaidsModel, error) {
	raidPartySize := len(raid.AvatarsOnRaids())

	if raidPartySize >= raid.PlayerLimit {
		return prisma.AvatarsOnRaidsModel{}, fmt.Errorf("party has reached its limit of %d", raidPartySize)
	}

	incrementedRaid, err := IncrementPartySize(db, raid.ID)
	if err != nil {
		return prisma.AvatarsOnRaidsModel{}, err
	}

	newRaidMemberPosition := incrementedRaid.PlayerCount

	raidMember, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Position.Set(newRaidMemberPosition),
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar.ID),
		),
	).Exec(context.Background())
	if err != nil {
		return raidMember, err
	}

	logger.GetLogger().Info("successfully added avatar to raid", zap.Int("avatarID", avatar.ID), zap.Int("raidID", raid.ID))

	return raidMember, nil
}

// IncrementPartySize increases the raid party size by one
func IncrementPartySize(db *prisma.PrismaClient, raidID int) (prisma.RaidModel, error) {
	// Prisma doesn't support incrementing values within one query yet. So we have to fetch twice. Once to get the original player count. Second to update that value.
	raid, err := db.Raid.FindOne(
		prisma.Raid.ID.Equals(raidID),
	).Exec(context.Background())
	if err != nil {
		return prisma.RaidModel{}, err
	}

	raidPartySize := len(raid.AvatarsOnRaids())

	if raidPartySize >= raid.PlayerLimit {
		return prisma.RaidModel{}, fmt.Errorf("party has reached its player limit of %d", raidPartySize)
	}

	incrementedRaid, err := db.Raid.FindOne(
		prisma.Raid.ID.Equals(raidID),
	).Update(
		prisma.Raid.PlayerCount.Set(raid.PlayerCount + 1)).Exec(context.Background())
	if err != nil {
		return prisma.RaidModel{}, err
	}

	return incrementedRaid, nil
}
