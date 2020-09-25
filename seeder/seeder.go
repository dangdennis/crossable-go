package seeder

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v5"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/repositories/users"
)

// Run runs the seeder
func Run() {
	db := prisma.Client()
	rando, _ := rand.Prime(rand.Reader, 128)

	// USERS
	gofakeit.Seed(rando.Int64())
	user1, err := users.CreateUser(db, users.UserAttrs{
		DiscordUserID:   strconv.FormatUint(uint64(gofakeit.Number(10000000, 90000000)), 10),
		Email:           toPtrString(gofakeit.Email()),
		DiscordUsername: toPtrString(gofakeit.Username()),
		FirstName:       toPtrString(gofakeit.FirstName()),
		LastName:        toPtrString(gofakeit.LastName()),
	})
	handleError(err)

	user2, err := users.CreateUser(db, users.UserAttrs{
		DiscordUserID:   strconv.FormatUint(uint64(gofakeit.Number(10000000, 90000000)), 10),
		Email:           toPtrString(gofakeit.Email()),
		DiscordUsername: toPtrString(gofakeit.Username()),
		FirstName:       toPtrString(gofakeit.FirstName()),
		LastName:        toPtrString(gofakeit.LastName()),
	})
	handleError(err)

	// AVATARS
	avatar1, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user1.ID),
		),
	).Exec(context.Background())
	handleError(err)

	avatar2, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user2.ID),
		),
	).Exec(context.Background())
	handleError(err)

	// RAID BOSSES
	bossLichKing, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("Arthas Menethil, The Lich King"),
		prisma.RaidBoss.Image.Set("https://cdn.vox-cdn.com/thumbor/k6m7tw54mdYa2yJoYbk3FuIYFZg=/0x0:1024x576/1920x0/filters:focal(0x0:1024x576):format(webp):no_upscale()/cdn.vox-cdn.com/uploads/chorus_asset/file/19748343/155054_the_lich_king.jpg"),
	).Exec(context.Background())
	handleError(err)

	bossAlienQueen, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("The Alien Queen"),
		prisma.RaidBoss.Image.Set("https://vignette.wikia.nocookie.net/avp/images/7/74/Promo07.PNG/revision/latest?cb=20120527102557"),
	).Exec(context.Background())
	handleError(err)

	now := time.Now()
	later := time.Now().Add(49 * time.Hour)

	// RAIDS
	raid1, err := db.Raid.CreateOne(
		prisma.Raid.Active.Set(false),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now.Add(1*time.Second)),
		prisma.Raid.EndTime.Set(later),
		prisma.Raid.CompletionProgress.Set(1.0),
	).Exec(context.Background())
	handleError(err)

	raid2, err := db.Raid.CreateOne(
		prisma.Raid.Active.Set(true),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now),
		prisma.Raid.EndTime.Set(later),
	).Exec(context.Background())
	handleError(err)

	// RAID BOSSES IN RAID
	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossLichKing.ID)),
	).Exec(context.Background())
	handleError(err)

	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid2.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossAlienQueen.ID)),
	).Exec(context.Background())
	handleError(err)

	// AVATARS ON RAIDS
	avatar1OnRaid1, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar1.ID),
		),
	).Exec(context.Background())
	handleError(err)

	avatar1OnRaid2, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid2.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar1.ID),
		),
	).Exec(context.Background())
	handleError(err)

	avatar2OnRaid1, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar2.ID),
		),
	).Exec(context.Background())
	handleError(err)

	avatar2OnRaid2, err := db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid2.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar2.ID),
		),
	).Exec(context.Background())
	handleError(err)

	fmt.Println(avatar1OnRaid1)
	fmt.Println(avatar1OnRaid2)
	fmt.Println(avatar2OnRaid1)
	fmt.Println(avatar2OnRaid2)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func toPtrString(str string) *string {
	return &str
}
