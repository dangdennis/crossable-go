package raid

import (
	"errors"
	"fmt"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/raids"
	"github.com/dangdennis/crossing/common/repositories/users"
)

// ErrExistingRaidMembership is an error for avatars attempting to join a raid they're already members of.
var ErrExistingRaidMembership error = errors.New("avatar is already a member of the raid")

// AssignAvatarToRaid assigns an avatar to the most recently started active raid
func AssignAvatarToRaid(discordID string) error {
	db := prisma.Client()

	raid, err := raids.FindLatestActiveRaid(db)
	if err != nil {
		return err
	}

	user, err := users.FindUserByDiscordID(db, discordID)
	if err != nil {
		return err
	}

	avatar, ok := user.Avatar()
	if !ok {
		return fmt.Errorf("user does not have an avatar")
	}

	raidMember, err := raids.GetAvatarRaidMembership(db, avatar, raid)
	if err == nil && raidMember.AvatarID > 0 {
		return ErrExistingRaidMembership
	}

	_, err = raids.JoinRaid(db, raid, avatar)
	if err != nil {
		return nil
	}

	return nil
}

