package repositories

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/generator/runtime"

	prisma "github.com/dangdennis/crossing/db"
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
		prisma.Raid.StartTime.Order(runtime.DESC),
	).Take(1).Exec(context.Background())
	if err != nil {
		return r, err
	}

	if len(raids) == 0 {
		return r, fmt.Errorf("failed to find the most recently started active raid")
	}

	return raids[0], nil
}
