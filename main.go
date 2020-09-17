package main

import (
	"github.com/dangdennis/crossing/db"
)

func main() {
	client := db.NewClient()
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := client.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	// ctx := context.Background()

	// // create a user
	// createdUser, err := client.User.CreateOne(
	// 	db.User.Email.Set("john.doe@example.com"),
	// 	db.User.Name.Set("John Doe"),
	// 	db.User.Age.Set(5),

	// 	// ID is optional, which is why it's specified last. if you don't set it
	// 	// an ID is auto generated for you
	// 	db.User.ID.Set("123"),
	// ).Exec(ctx)

	// log.Printf("created user: %+v", createdUser)

	// // find a single user
	// user, err := client.User.FindOne(
	// 	db.User.Email.Equals("john.doe@example.com"),
	// ).Exec(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Printf("user: %+v", user)

	// // for optional/nullable values, you need to check the function and create two return values
	// // `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
	// // `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
	// // `ok` is true and `name` will be "John Doe".
	// name, ok := user.Name()

	// if !ok {
	// 	log.Printf("user's name is null")
	// 	return
	// }

	// log.Printf("The users's name is: %s", name)
}

// Variables used for command line parameters
// var (
// 	Token string
// )

// func init() {
// 	flag.StringVar(&Token, "t", "", "Bot Token")
// 	flag.Parse()
// }

// func main() {

// 	// Create a new Discord session using the provided bot token.
// 	dg, err := discordgo.New("Bot " + Token)
// 	if err != nil {
// 		fmt.Println("error creating Discord session,", err)
// 		return
// 	}

// 	// Register the messageCreate func as a callback for MessageCreate events.
// 	dg.AddHandler(messageCreate)

// 	// In this example, we only care about receiving message events.
// 	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

// 	// Open a websocket connection to Discord and begin listening.
// 	err = dg.Open()
// 	if err != nil {
// 		fmt.Println("error opening connection,", err)
// 		return
// 	}

// 	// Wait here until CTRL-C or other term signal is received.
// 	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
// 	sc := make(chan os.Signal, 1)
// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
// 	<-sc

// 	// Cleanly close down the Discord session.
// 	dg.Close()
// }

// // This function will be called (due to AddHandler above) every time a new
// // message is created on any channel that the authenticated bot has access to.
// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	fmt.Println(s)
// 	fmt.Println(m)
// 	// Ignore all messages created by the bot itself
// 	// This isn't required in this specific example but it's a good practice.
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}
// 	// If the message is "ping" reply with "Pong!"
// 	if m.Content == "ping" {
// 		s.ChannelMessageSend(m.ChannelID, "Pong!")
// 	}

// 	// If the message is "pong" reply with "Ping!"
// 	if m.Content == "pong" {
// 		s.ChannelMessageSend(m.ChannelID, "Ping!")
// 	}
// }
