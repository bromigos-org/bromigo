package commands

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const ventChannelName = "vent-anonymously" // Replace with your actual channel name
func VentAnonymously(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Find the channel (or thread) where the message was posted
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Error retrieving channel: %v\n", err)
		return
	}

	// Check if the message was posted in a thread under the vent-anonymously channel
	isThread := channel.Type == discordgo.ChannelTypeGuildPublicThread || channel.Type == discordgo.ChannelTypeGuildPrivateThread
	if isThread {
		parentChannel, err := s.Channel(channel.ParentID)
		if err != nil {
			log.Printf("Error retrieving parent channel: %v\n", err)
			return
		}

		if parentChannel.Name == ventChannelName {
			// Delete the thread
			_, err := s.ChannelDelete(m.ChannelID)
			if err != nil {
				log.Printf("Error deleting thread: %v\n", err)
				return
			}

			threadStart := discordgo.ThreadStart{
				Name: channel.Name,
			}

			// Define thread creation parameters
			threadCreateData := discordgo.MessageSend{
				// Content: m.Content,
				Content: "This thread was created to replace a deleted thread.",
			}

			// Create a new thread with the same name as the deleted thread
			newThread, err := s.ForumThreadStartComplex(parentChannel.ID, &threadStart, &threadCreateData)
			if err != nil {
				log.Printf("Error creating new thread: %v\n", err)
				return
			}
			time.Sleep(2 * time.Second)
			// Post the message content as the bot in the new thread
			_, err = s.ChannelMessageSend(newThread.ID, m.ContentWithMentionsReplaced())
			if err != nil {
				log.Printf("Error sending message in new thread: %v\n", err)
			}
		}
	}
}
