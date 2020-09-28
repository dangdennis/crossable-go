package stories

import (
	"context"

	prisma "github.com/dangdennis/crossing/db"
)

// CreateStory creates a new stories
func CreateStory(db *prisma.PrismaClient) (prisma.StoryModel, error) {
	return db.Story.CreateOne().Exec(context.Background())
}
