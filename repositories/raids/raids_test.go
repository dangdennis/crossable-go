package raids

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dangdennis/crossing/db"
)

func TestFindWeeklyActiveRaid(t *testing.T) {
	// TODO seed data within the test
	raid, err := FindWeeklyActiveRaid(db.Client())
	require.NoError(t, err)

	require.True(t, raid.ID > 0)
}

func TestJoinRaid(t *testing.T) {

}
