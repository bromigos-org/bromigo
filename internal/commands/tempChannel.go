package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ChannelInfo stores if a channel was created by the bot and its associated text channel ID
type ChannelInfo struct {
	CreatedByBot  bool
	TextChannelID string // ID of the associated text chat
}

var (
	botCreatedChannels = make(map[string]ChannelInfo) // TODO: Use a database to store this information instead of an in-memory map.
	mutex              sync.Mutex                     // Mutex to protect access to botCreatedChannels
)

// VoiceStateUpdate handles voice state changes, creating or deleting temporary voice channels
func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	const triggerChannelName = "=Join To Start Game" // Name of the trigger channel

	// Find the trigger channel ID dynamically by name
	channels, err := s.GuildChannels(v.GuildID)
	if err != nil {
		fmt.Printf("Error retrieving channels: %v\n", err)
		return
	}

	var triggerChannel *discordgo.Channel
	for _, channel := range channels {
		if channel.Name == triggerChannelName && channel.Type == discordgo.ChannelTypeGuildVoice {
			triggerChannel = channel
			break
		}
	}

	if triggerChannel == nil {
		fmt.Println("Trigger channel not found")
		return
	}
	// User joins the trigger channel
	if (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID != triggerChannel.ID) && v.ChannelID == triggerChannel.ID {
		time.Sleep(1 * time.Second) // Delay to ensure any previous channel deletion is processed

		// Create a new voice channel with a unique name
		// Attempt to get the member from the state
		member, err := s.State.Member(v.GuildID, v.UserID)
		if err != nil { // Member not found in state, fall back to API call
			member, err = s.GuildMember(v.GuildID, v.UserID)
			if err != nil {
				fmt.Printf("Error retrieving guild member: %v\n", err)
				return
			}
		}

		// Use member's nickname if available, otherwise use user's username
		nameToUse := member.User.GlobalName
		if member.Nick != "" {
			nameToUse = member.Nick
		}

		channelName := fmt.Sprintf("%s's Channel", nameToUse)
		channel, err := s.GuildChannelCreateComplex(v.GuildID, discordgo.GuildChannelCreateData{
			Name:     channelName,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: triggerChannel.ParentID, // Set the new channel in the same category as the trigger channel
		})
		if err != nil {
			fmt.Printf("Error creating voice channel: %v\n", err)
			return
		}

		mutex.Lock()
		// Assuming "TextChannelID" is the ID of the text channel you want to associate with this voice channel
		botCreatedChannels[channel.ID] = ChannelInfo{CreatedByBot: true, TextChannelID: channel.ID}
		mutex.Unlock()

		// Move the user to the newly created voice channel
		err = s.GuildMemberMove(v.GuildID, v.UserID, &channel.ID)
		if err != nil {
			fmt.Printf("Error moving user to new channel: %v\n", err)
			s.ChannelDelete(channel.ID) // Delete the channel if the user couldn't be moved
			return
		}
	} else if v.BeforeUpdate != nil && v.ChannelID != triggerChannel.ID {
		// User leaves a bot-created channel
		mutex.Lock()
		channelInfo, exists := botCreatedChannels[v.BeforeUpdate.ChannelID]
		mutex.Unlock()

		if exists && channelInfo.CreatedByBot {
			guild, err := s.State.Guild(v.GuildID)
			if err != nil {
				fmt.Printf("Error retrieving guild: %v\n", err)
				return
			}

			// Check if the channel is empty
			isEmpty := true
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == v.BeforeUpdate.ChannelID {
					isEmpty = false
					break
				}
			}

			// Delete the channel if empty
			if isEmpty {
				_, err := s.ChannelDelete(v.BeforeUpdate.ChannelID)
				if err != nil {
					return
				}

				mutex.Lock()
				delete(botCreatedChannels, v.BeforeUpdate.ChannelID)
				mutex.Unlock()
			}
		}
	}
}

// MessageReactionAdd handles reaction additions, renaming channels based on reactions
func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	mutex.Lock()
	defer mutex.Unlock()

	// Iterate over botCreatedChannels to find a match for TextChannelID
	for voiceChannelID, info := range botCreatedChannels {
		if info.TextChannelID == m.ChannelID && m.Emoji.Name == "bromigo" { // Use the emoji you want to trigger the rename
			// Retrieve the message to get its content for the new channel name
			msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
			if err != nil {
				fmt.Printf("Error retrieving message: %v\n", err)
				return
			}

			// Use the message content as the new channel name, ensure it meets Discord's channel naming requirements
			newName := msg.Content
			if len(newName) > 100 { // Discord's maximum channel name length is 100 characters
				newName = newName[:100]
			}

			channelConfig := &discordgo.ChannelEdit{
				Name: newName,
			}

			// Edit the voice channel with the new name derived from the message content
			_, err = s.ChannelEdit(voiceChannelID, channelConfig)
			if err != nil {
				fmt.Printf("Error renaming channel: %v\n", err)
			}
			break // Exit the loop after renaming
		}
	}
}
