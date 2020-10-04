package users

import (
	"context"
	"fmt"

	prisma "github.com/dangdennis/crossing/common/db"
)

// FindUserByDiscordID finds a user entity by their discord id and their avatar
func FindUserByDiscordID(db *prisma.PrismaClient, discordID string) (prisma.UserModel, error) {
	return db.User.FindOne(
		prisma.User.DiscordUserID.Equals(discordID),
	).With(
		prisma.User.Avatar.Fetch(),
		prisma.User.Wallet.Fetch(),
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
func CreateAvatar(db *prisma.PrismaClient, userID int) (prisma.AvatarModel, error) {
	avatar, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(userID),
		)).Exec(context.Background())
	if err != nil {
		return avatar, err
	}

	return avatar, nil
}

// CreateWallet creates an avatar
func CreateWallet(db *prisma.PrismaClient, userID int) (prisma.WalletModel, error) {
	avatar, err := db.Wallet.CreateOne(
		prisma.Wallet.User.Link(
			prisma.User.ID.Equals(userID),
		)).Exec(context.Background())
	if err != nil {
		return avatar, err
	}

	return avatar, nil
}

// AwardTokens awards an amount of tokens to a user
func AwardTokens(db *prisma.PrismaClient, userID int, amount int) error {
	user, err := db.User.FindOne(
		prisma.User.ID.Equals(userID),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to find user. err=%w", err)
	}

	wallet, ok := user.Wallet()
	if !ok {
		return fmt.Errorf("failed to find user wallet")
	}

	_, err = db.Wallet.FindOne(
		prisma.Wallet.ID.Equals(wallet.ID),
	).Update(
		prisma.Wallet.Balance.Set(wallet.Balance + amount),
	).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to update wallet. err=%w", err)
	}

	return nil
}
