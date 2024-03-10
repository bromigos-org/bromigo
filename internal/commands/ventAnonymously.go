package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bromigos-org/bromigo/internal/utils"
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

	utils.LogJSON("Channel", channel)

	utils.LogJSON("Message", m)

	// Check if the message was posted in a thread under the vent-anonymously channel
	isThread := channel.Type == discordgo.ChannelTypeGuildPublicThread || channel.Type == discordgo.ChannelTypeGuildPrivateThread
	if isThread {
		parentChannel, err := s.Channel(channel.ParentID)
		if err != nil {
			log.Printf("Error retrieving parent channel: %v\n", err)
			return
		}

		if parentChannel.Name == ventChannelName {

			fmt.Println(channel.Name)

			threadStart := discordgo.ThreadStart{
				Name: channel.Name,
			}

			mess, err := s.ChannelMessage(channel.ID, channel.LastMessageID)
			if err != nil {
				log.Printf("Error retrieving message: %v\n", err)
				return
			}

			fmt.Printf("Message: %+v\n", mess)

			// // Get first message
			// messages, err := s.ChannelMessages(channel.ID, 1, "", "", "")
			// if err != nil {
			// 	log.Printf("Error retrieving messages: %v\n", err)
			// 	return
			// }

			// // Check if there are any messages to avoid index out-of-range errors
			// if len(messages) > 0 {
			// 	// Print the first message content. No need to use & as we're accessing a field of the struct pointed to by the slice element.
			// 	fmt.Printf("First Message Content: %s\n", messages[0].Content)

			// 	// If you want to print the entire message struct, use %+v to print with field names for readability.
			// 	fmt.Printf("First Message Struct: %+v\n", *messages[0])
			// } else {
			// 	fmt.Println("No messages found in the thread.")
			// }

			// Define thread creation parameter bs
			threadCreateData := discordgo.MessageSend{
				// Content: m.Content,
				Content: "This thread was created to replace a deleted thread.",
			}

			// // Delete the thread
			// _, err = s.ChannelDelete(channel.ID)
			// if err != nil {
			// 	log.Printf("Error deleting thread: %v\n", err)
			// 	return
			// }

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
