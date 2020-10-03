package stories

import (
	"context"
	"fmt"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/messages"
)

// CreateStory creates a new stories
func CreateStory(db *prisma.PrismaClient) (prisma.StoryModel, error) {
	return db.Story.CreateOne().Exec(context.Background())
}

// GetCurrentEventInStory gets the next event in the story sequence that hasn't occurred yet
func GetCurrentEventInStory(db *prisma.PrismaClient, story prisma.StoryModel) (prisma.EventModel, error) {
	events, err := db.Event.FindMany(
		prisma.Event.StoryID.Equals(story.ID),
		prisma.Event.Occurred.Equals(false),
	).OrderBy(
		prisma.Event.Sequence.Order(prisma.ASC),
	).Take(1).Exec(context.Background())

	if err != nil {
		return prisma.EventModel{}, fmt.Errorf("failed to find events for story. err=%w", err)
	}

	if len(events) == 0 {
		return prisma.EventModel{}, fmt.Errorf("no active event found for current storyline. err=%w", err)
	}

	return events[0], nil
}

// GetEventIntroMessage gets the intro story for a particular event
func GetEventIntroMessage(db *prisma.PrismaClient, event prisma.EventModel) (prisma.MessageModel, error) {
	introMessages, err := db.Message.FindMany(
		prisma.Message.EventID.Equals(event.ID),
		prisma.Message.Type.Equals(messages.MessageTypeEventIntro.String()),
	).Take(1).Exec(context.Background())
	if err != nil {
		return prisma.MessageModel{}, fmt.Errorf("failed to find an intro message for the event. err=%w", err)
	}

	return introMessages[0], nil
}

// GetEventOutroMessage gets the intro story for a particular event
func GetEventOutroMessage(db *prisma.PrismaClient, event prisma.EventModel) (prisma.MessageModel, error) {
	introMessages, err := db.Message.FindMany(
		prisma.Message.EventID.Equals(event.ID),
		prisma.Message.Type.Equals(messages.MessageTypeEventOutro.String()),
	).Take(1).Exec(context.Background())
	if err != nil {
		return prisma.MessageModel{}, fmt.Errorf("failed to find an intro message for the event. err=%w", err)
	}

	return introMessages[0], nil
}

// GetActionMessageForEventAndRaidMember find the relevant player message for the avatar.
// If an action message does not exist for the current raid member, use the last one in the message sequence as the default.
func GetActionMessageForEventAndRaidMember(db *prisma.PrismaClient, event prisma.EventModel, raidMember prisma.AvatarsOnRaidsModel) (prisma.MessageModel, error) {
	manyMessages, err := db.Message.FindMany(
		prisma.Message.EventID.Equals(event.ID),
		prisma.Message.Type.Equals(messages.MessageTypeActionSingle.String()),
		prisma.Message.Sequence.Equals(raidMember.Position),
		prisma.Message.Default.Equals(false),
	).Take(1).Exec(context.Background())
	if err != nil {
		defaultMessages, err := db.Message.FindMany(
			prisma.Message.EventID.Equals(event.ID),
			prisma.Message.Type.Equals(messages.MessageTypeActionSingle.String()),
			prisma.Message.Default.Equals(true),
		).Take(1).Exec(context.Background())
		if err != nil {
			return prisma.MessageModel{}, fmt.Errorf("failed to find a message for the action. err=%w", err)
		}

		if len(defaultMessages) == 0 {
			return prisma.MessageModel{}, fmt.Errorf("failed to find a default message for avatar action")
		}

		return defaultMessages[0], nil
	}

	if len(manyMessages) == 0 {
		return prisma.MessageModel{}, fmt.Errorf("failed to find a message for avatar action")
	}

	return manyMessages[0], nil
}

// CreateAvatarEventAction logs the avatar's action for the event
func CreateAvatarEventAction(db *prisma.PrismaClient, currentEvent prisma.EventModel, avatar prisma.AvatarModel) error {
	_, err := db.Action.CreateOne(
		prisma.Action.Event.Link(
			prisma.Event.ID.Equals(currentEvent.ID),
		),
		prisma.Action.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar.ID)),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create action for avatar. err=%w", err)
	}

	return nil
}
