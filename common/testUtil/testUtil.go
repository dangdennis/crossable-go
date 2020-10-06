package testUtil

import (
	"context"
	"time"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/messages"
)

type Mocks struct {
	Raid    *prisma.RaidModel
	Message *prisma.MessageModel
	Event   *prisma.EventModel
	Story   *prisma.StoryModel
}

// NewMocks initializes some mock data into db, and stored into memory
func NewMocks(db *prisma.PrismaClient) (*Mocks, error) {
	t := &Mocks{}

	mockStory, err := db.Story.CreateOne().Exec(context.Background())
	if err != nil {
		return t, err
	}
	t.Story = &mockStory

	mockEvent1, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(mockStory.ID),
		),
		prisma.Event.Sequence.Set(1),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	t.Event = &mockEvent1

	content := "some random message"

	mockMessage, err := db.Message.CreateOne(
		prisma.Message.Event.Link(
			prisma.Event.ID.Equals(mockEvent1.ID)),
		prisma.Message.Content.Set(content),
		prisma.Message.Type.Set(messages.MessageTypeEventIntro.String()),
		prisma.Message.Sequence.Set(1),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	t.Message = &mockMessage

	mockRaid, err := db.Raid.CreateOne(
		prisma.Raid.Story.Link(
			prisma.Story.ID.Equals(mockStory.ID)),
		// Set the start date super far in the future to ensure we query for it every time
		prisma.Raid.StartTime.Set(time.Now().AddDate(10, 0, 0)),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	t.Raid = &mockRaid

	return t, nil
}

// Cleans up the mocks created from the db
func (t *Mocks) Cleanup(db *prisma.PrismaClient) error {
	if t.Message != nil {
		_, err := db.Message.FindOne(
			prisma.Message.ID.Equals(t.Message.ID),
		).Delete().Exec(context.Background())
		if err != nil {
			return err
		}
	}

	if t.Event != nil {
		_, err := db.Event.FindOne(
			prisma.Event.ID.Equals(t.Event.ID),
		).Delete().Exec(context.Background())
		if err != nil {
			return err
		}
	}

	if t.Raid != nil {
		_, err := db.Raid.FindOne(
			prisma.Raid.ID.Equals(t.Raid.ID),
		).Delete().Exec(context.Background())
		if err != nil {
			return err
		}
	}

	if t.Story != nil {
		_, err := db.Story.FindOne(
			prisma.Story.ID.Equals(t.Story.ID),
		).Delete().Exec(context.Background())
		if err != nil {
			return err
		}
	}



	return nil
}

