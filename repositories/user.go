package repositories

import (
	"context"
	"fmt"

	prisma "github.com/dangdennis/crossing/db"
)

// FindUserByDiscordID finds a user entity by their discord id
func FindUserByDiscordID(db *prisma.PrismaClient, discordID string) (prisma.UserModel, error) {
	return db.User.FindOne(
		prisma.User.DiscordUserID.Equals(discordID),
	).Exec(context.Background())
}

// UserAttrs defines the request `create` payload
type UserAttrs struct {
	DiscordUserID   string
	Email           *string
	DiscordUsername *string
	FirstName       *string
	LastName        *string
}

// CreateUser creates a user
func CreateUser(db *prisma.PrismaClient, attrs UserAttrs) (prisma.UserModel, error) {
	fmt.Println("creating a user")

	user, err := db.User.CreateOne(
		prisma.User.DiscordUserID.Set(attrs.DiscordUserID),
		prisma.User.DiscordUsername.SetOptional(attrs.DiscordUsername),
		prisma.User.Email.SetOptional(attrs.Email),
		prisma.User.FirstName.SetOptional(attrs.FirstName),
		prisma.User.LastName.SetOptional(attrs.LastName),
	).Exec(context.Background())
	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateAvatar creates an avatar
func CreateAvatar(db *prisma.PrismaClient, userID uint64) (prisma.AvatarModel, error) {
	fmt.Println("creating an avatar")

	avatar, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(int(userID)),
		)).Exec(context.Background())
	if err != nil {
		return avatar, err
	}

	return avatar, nil
}
