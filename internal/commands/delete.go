package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DeleteMessages deletes the last N messages in the current channel for users with Administrator permission
func DeleteMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Split the command to get the number of messages to delete
	parts := strings.Fields(m.Content)
	if len(parts) != 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: @bot mdelete <number>")
		return // Incorrect command usage
	}

	// Parse the number of messages to delete
	numMessages, err := strconv.Atoi(parts[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of messages to delete.")
		return // Not a valid number
	}

	// Check if the user has Administrator permissions
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		fmt.Printf("Error retrieving guild member: %v\n", err)
		return
	}

	hasAdminPermission := false
	for _, roleID := range member.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		if err != nil {
			fmt.Printf("Error retrieving role: %v\n", err)
			continue // Skip on error
		}

		// Check if the role has Administrator permission
		// Administrator permission is represented by the bit flag discordgo.PermissionAdministrator
		// If the bitwise AND operation between the role's permissions and discordgo.PermissionAdministrator is not 0, the role has Administrator permission
		// This is because the Administrator permission is the highest bit in the permission integer
		// If the bit is set, the value is a power of 2, which results in a non-zero value when ANDed with any other number
		// For example, 8 (Administrator permission) AND 8 is 8, which is non-zero
		// 8 AND 16 is 0, as the bits do not overlap
		// https://github.com/discord/discord-api-docs/blob/main/docs/topics/Permissions.md
		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			hasAdminPermission = true
			break
		}
	}

	if !hasAdminPermission {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return // User does not have Administrator permissions
	}

	// Retrieve the last N messages in the channel
	messages, err := s.ChannelMessages(m.ChannelID, numMessages, "", "", "")
	if err != nil {
		fmt.Printf("Error retrieving messages: %v\n", err)
		return
	}

	// Delete the messages
	for _, message := range messages {
		err := s.ChannelMessageDelete(m.ChannelID, message.ID)
		if err != nil {
			fmt.Printf("Error deleting message: %v\n", err)
			// Continue attempting to delete other messages even if one fails
		}
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleted the last %d messages.", numMessages))
}
