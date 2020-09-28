package story

import (
	"context"

	prisma "github.com/dangdennis/crossing/db"
)

// CreateStory creates a new story
func CreateStory(db *prisma.PrismaClient) (prisma.StoryModel, error) {
	return db.Story.CreateOne().Exec(context.Background())
}
