package raid

import (
	"testing"

	"github.com/stretchr/testify/require"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/testUtil"
)


func TestAssignAvatarToRaid(t *testing.T) {
	db := prisma.Client()
	tu, err := testUtil.NewMocks(db)
	require.NoError(t, err)

	err = tu.Cleanup(db)
	require.NoError(t, err)
}
