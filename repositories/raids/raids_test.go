package raids

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dangdennis/crossing/db"
)

func TestFindWeeklyActiveRaid(t *testing.T) {
	raid, err := FindWeeklyActiveRaid(db.Client())
	require.NoError(t, err)

	require.True(t, raid.ID > 0)
}
