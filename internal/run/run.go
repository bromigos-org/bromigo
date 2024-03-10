package run

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bromigos-org/bromigo/internal/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load() // Load .env file if it exists

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Println("Warning: DISCORD_BOT_TOKEN not set in .env, ensure it's set in your environment")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("Error creating Discord session: %v\n", err)
		return
	}

	// Register event handlers
	dg.AddHandler(onReady)
	dg.AddHandler(onDisconnect)
	dg.AddHandler(onReconnect)

	dg.AddHandler(commands.BotMention)
	dg.AddHandler(commands.VentAnonymously)
	dg.AddHandler(commands.HandleThreadMessages)
	dg.AddHandler(commands.VoiceStateUpdate)
	dg.AddHandler(commands.MessageReactionAdd)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMessageReactions | discordgo.IntentMessageContent

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Printf("Error opening connection: %v\n", err)
		return
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close() // Cleanly close down the Discord session.
}

// onReady is called when the bot is ready to start receiving events.
func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Bot is ready.")
}

// onDisconnect is called when the bot disconnects from Discord.
func onDisconnect(s *discordgo.Session, event *discordgo.Disconnect) {
	log.Println("Bot disconnected.")
}

// onReconnect is called when the bot reconnects to Discord.
func onReconnect(s *discordgo.Session, event *discordgo.Connect) {
	log.Println("Bot reconnected.")
}
