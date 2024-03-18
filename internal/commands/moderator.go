package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const moderatorChannelName = "moderator-only"

// PostToChannel handles messages in #moderator-only to post content to a specified channel
func PostToChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf("Error retrieving channel: %v\n", err)
		return
	}

	// Check if the message was posted in the moderator-only channel
	if channel.Name == moderatorChannelName {
		// Extract channel ID from the mention format
		rgx := regexp.MustCompile(`<#(\d+)>`)
		matches := rgx.FindStringSubmatch(m.Content)
		if len(matches) < 2 {
			fmt.Println("Invalid channel mention format")
			return
		}
		channelID := matches[1]

		// Find the end of the channel mention
		mentionEndIndex := strings.Index(m.Content, fmt.Sprintf("<#%s>", channelID)) + len(fmt.Sprintf("<#%s>", channelID))
		if mentionEndIndex == -1 {
			fmt.Println("Channel mention not found")
			return
		}

		// Extract everything after the channel mention as the message content
		messageContent := strings.TrimSpace(m.Content[mentionEndIndex:])
		if messageContent == "" {
			fmt.Println("No message content provided")
			return
		}

		// Send the message to the specified channel
		_, err := s.ChannelMessageSend(channelID, messageContent)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			return
		}
	}
}
