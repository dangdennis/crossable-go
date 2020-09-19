package repositories

import (
	"context"
	"fmt"

	"github.com/dangdennis/crossing/db"
)

// FindUserByDiscordID finds a user entity by their discord id
func FindUserByDiscordID(client *db.PrismaClient, discordID string) (db.UserModel, error) {
	return client.User.FindOne(
		db.User.ID.Equals(discordID),
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
func CreateUser(client *db.PrismaClient, attrs UserAttrs) (db.UserModel, error) {
	fmt.Println("creating a user")

	user, err := client.User.CreateOne(
		db.User.DiscordUserID.Set(attrs.DiscordUserID),
		db.User.DiscordUsername.SetOptional(attrs.DiscordUsername),
		db.User.Email.SetOptional(attrs.Email),
		db.User.FirstName.SetOptional(attrs.FirstName),
		db.User.LastName.SetOptional(attrs.LastName),
	).Exec(context.Background())
	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateAvatar creates an avatar
func CreateAvatar(client *db.PrismaClient, userID string) (db.AvatarModel, error) {
	fmt.Println("creating an avatar")

	avatar, err := client.Avatar.CreateOne(
		db.Avatar.User.Link(
			db.User.ID.Equals(userID),
		)).Exec(context.Background())
	if err != nil {
		return avatar, err
	}

	return avatar, nil
}
