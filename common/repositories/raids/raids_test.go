package raids

import (
	"crypto/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/require"

	"github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/users"
)

func TestCreateRaid(t *testing.T) {
	raid, err := CreateRaid(db.Client())
	require.NoError(t, err)
	require.True(t, raid.ID > 0)
}

func TestFindWeeklyActiveRaid(t *testing.T) {
	// TODO seed data within the test
	raid, err := FindLatestActiveRaid(db.Client())
	require.NoError(t, err)
	require.True(t, raid.ID > 0)
}

func TestJoinRaid(t *testing.T) {
	client := db.Client()
	prime, _ := rand.Prime(rand.Reader, 128)
	gofakeit.Seed(prime.Int64())

	user, err := users.CreateUser(client, users.UserAttrs{
		DiscordUserID: gofakeit.UUID(),
	})
	require.NoError(t, err)

	avatar, err := users.CreateAvatar(client, user.ID)
	require.NoError(t, err)

	raid, err := CreateRaid(client)
	require.NoError(t, err)

	raidMember, err := JoinRaid(client, raid, avatar)
	require.NoError(t, err)
	require.True(t, raidMember.AvatarID == avatar.ID)
	require.True(t, raidMember.RaidID == raid.ID)
}
