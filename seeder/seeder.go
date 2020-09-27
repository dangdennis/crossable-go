package seeder

import (
	"context"
	"crypto/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v5"

	prisma "github.com/dangdennis/crossing/db"
	"github.com/dangdennis/crossing/repositories/users"
)

// Run runs the seeder
func Run() {
	db := prisma.Client()
	random, _ := rand.Prime(rand.Reader, 128)
	ctx := context.Background()

	// Create a couple users
	gofakeit.Seed(random.Int64())
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

	// Create the avatars for the new users
	avatar1, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user1.ID),
		),
	).Exec(ctx)
	handleError(err)

	avatar2, err := db.Avatar.CreateOne(
		prisma.Avatar.User.Link(
			prisma.User.ID.Equals(user2.ID),
		),
	).Exec(ctx)
	handleError(err)

	// Create a couple starter raid bosses
	bossLichKing, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("Arthas Menethil, The Lich King"),
		prisma.RaidBoss.Image.Set("https://cdn.vox-cdn.com/thumbor/k6m7tw54mdYa2yJoYbk3FuIYFZg=/0x0:1024x576/1920x0/filters:focal(0x0:1024x576):format(webp):no_upscale()/cdn.vox-cdn.com/uploads/chorus_asset/file/19748343/155054_the_lich_king.jpg"),
	).Exec(ctx)
	handleError(err)

	bossAlienQueen, err := db.RaidBoss.CreateOne(
		prisma.RaidBoss.Name.Set("The Alien Queen"),
		prisma.RaidBoss.Image.Set("https://vignette.wikia.nocookie.net/avp/images/7/74/Promo07.PNG/revision/latest?cb=20120527102557"),
	).Exec(ctx)
	handleError(err)

	now := time.Now()
	later := time.Now().Add(49 * time.Hour)

	// Create a couple raids
	// The first is inactive. The second is active, open for another 7 days from time of seed.
	raid1, err := db.Raid.CreateOne(
		prisma.Raid.Active.Set(false),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now.Add(1*time.Second)),
		prisma.Raid.EndTime.Set(later),
		prisma.Raid.CompletionProgress.Set(1.0),
	).Exec(ctx)
	handleError(err)

	raid2, err := db.Raid.CreateOne(
		prisma.Raid.Active.Set(true),
		prisma.Raid.PlayerLimit.Set(20),
		prisma.Raid.StartTime.Set(now),
		prisma.Raid.EndTime.Set(later),
	).Exec(ctx)
	handleError(err)

	// Add the raid bosses to a raid
	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossLichKing.ID)),
	).Exec(ctx)
	handleError(err)

	_, err = db.RaidBossesOnRaids.CreateOne(
		prisma.RaidBossesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid2.ID)),
		prisma.RaidBossesOnRaids.RaidBoss.Link(
			prisma.RaidBoss.ID.Equals(bossAlienQueen.ID)),
	).Exec(ctx)
	handleError(err)

	// Add the avatars to the first raid.
	_, err = db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Position.Set(1),
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar1.ID),
		),
	).Exec(ctx)
	handleError(err)

	_, err = db.AvatarsOnRaids.CreateOne(
		prisma.AvatarsOnRaids.Position.Set(2),
		prisma.AvatarsOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(raid1.ID),
		),
		prisma.AvatarsOnRaids.Avatar.Link(
			prisma.Avatar.ID.Equals(avatar2.ID),
		),
	).Exec(ctx)
	handleError(err)

	// Create the Alien Queen story
	alienStory, err := db.Story.CreateOne().Exec(ctx)
	handleError(err)

	// Attach the alien story to the second created raid
	_, err = db.StoriesOnRaids.CreateOne(
		prisma.StoriesOnRaids.Raid.Link(
			prisma.Raid.ID.Equals(alienStory.ID),
		),
		prisma.StoriesOnRaids.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID),
		),
	).Exec(ctx)
	handleError(err)

	// Create the events for the story
	// Days 1 through 6
	alienStoryEvent1, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

	alienStoryEvent2, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

	alienStoryEvent3, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

	alienStoryEvent4, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

	alienStoryEvent5, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

	alienStoryEvent6, err := db.Event.CreateOne(
		prisma.Event.Story.Link(
			prisma.Story.ID.Equals(alienStory.ID)),
		prisma.Event.Sequence.Set(1),
	).Exec(ctx)

}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func toPtrString(str string) *string {
	return &str
}
