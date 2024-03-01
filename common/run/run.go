package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bromigos-org/bromigo/common/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load() // Load .env file if it exists

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Warning: DISCORD_BOT_TOKEN not set in .env, ensure it's set in your environment")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(commands.MessageCreate)

	// Register a new handler for voice state updates.
	// You will need to implement this function in your 'commands' package.
	dg.AddHandler(commands.VoiceStateUpdate)

	// Update intents to include voice states along with guild messages.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close() // Cleanly close down the Discord session.
}
