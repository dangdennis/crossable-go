package raid

import (
	"errors"
	"fmt"

	"github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/raids"
	"github.com/dangdennis/crossing/common/repositories/users"
)

// ErrExistingRaidMembership is an error for avatars attempting to join a raid they're already members of.
var ErrExistingRaidMembership error = errors.New("avatar is already a member of the raid")

// AssignAvatarToRaid assigns an avatar to the most recently started active raid
func AssignAvatarToRaid(client *db.PrismaClient, discordID string) error {
	raid, err := raids.FindLatestActiveRaid(client)
	if err != nil {
		return err
	}

	user, err := users.FindUserByDiscordID(client, discordID)
	if err != nil {
		return err
	}

	avatar, ok := user.Avatar()
	if !ok {
		return fmt.Errorf("user does not have an avatar")
	}

	raidMember, err := raids.GetAvatarRaidMembership(client, avatar, raid)
	if err == nil && raidMember.AvatarID > 0 {
		return ErrExistingRaidMembership
	}

	_, err = raids.JoinRaid(client, raid, avatar)
	if err != nil {
		return nil
	}

	return nil
}
