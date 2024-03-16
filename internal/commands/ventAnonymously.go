package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const ventChannelName = "vent-anonymously"

// Use a map to track the IDs of threads recreated by the bot for anonymity.
// TODO: Use a database to store this information instead of an in-memory map.
var anonymousThreads = make(map[string]bool)

// VentAnonymously is a handler for messages in the vent-anonymously channel
func VentAnonymously(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Error retrieving channel: %v\n", err)
		return
	}

	// Check if the message was posted in a thread under the vent-anonymously channel
	if isThreadUnderVentChannel(s, channel) {
		parentChannel, err := s.Channel(channel.ParentID)
		if err != nil {
			log.Printf("Error retrieving parent channel: %v\n", err)
			return
		}

		if _, managed := anonymousThreads[channel.ID]; !managed {
			// Delete the original thread
			_, err = s.ChannelDelete(channel.ID)
			if err != nil {
				log.Printf("Error deleting original thread: %v\n", err)
				return
			}

			threadStart := discordgo.ThreadStart{
				Name: channel.Name,
			}

			// Define thread creation parameters for the new anonymous thread
			threadCreateData := discordgo.MessageSend{
				Content: m.Content,
			}

			// Recreate the thread anonymously
			newThread, err := s.ForumThreadStartComplex(parentChannel.ID, &threadStart, &threadCreateData)
			if err != nil {
				log.Printf("Error creating new anonymous thread: %v\n", err)
				return
			}

			// Mark the new thread as managed
			anonymousThreads[newThread.ID] = true
		}
	}
}

// HandleThreadMessages is a handler for messages in threads under the vent-anonymously channel
func HandleThreadMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself or messages not in anonymous threads
	if m.Author.ID == s.State.User.ID || !anonymousThreads[m.ChannelID] {
		return
	}

	// Delete the original message
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Printf("Error deleting message: %v\n", err)
		return
	}

	// Repost the message content as the bot in the same thread
	_, err = s.ChannelMessageSend(m.ChannelID, "Anonymous: "+m.Content)
	if err != nil {
		log.Printf("Error reposting message as bot: %v\n", err)
	}
}

// isThreadUnderVentChannel checks if a thread is under the vent-anonymously channel
func isThreadUnderVentChannel(s *discordgo.Session, channel *discordgo.Channel) bool {
	if channel.Type != discordgo.ChannelTypeGuildPublicThread && channel.Type != discordgo.ChannelTypeGuildPrivateThread {
		return false
	}
	parentChannel, err := s.Channel(channel.ParentID)
	if err != nil {
		log.Printf("Error retrieving parent channel: %v\n", err)
		return false
	}
	return parentChannel.Name == ventChannelName
}
