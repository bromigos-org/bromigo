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
	botCreatedChannels = make(map[string]ChannelInfo)
	mutex              sync.Mutex // Mutex to protect access to botCreatedChannels
)

// VoiceStateUpdate handles voice state changes, creating or deleting temporary voice channels
func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	triggerChannelID := "1211873008806010930" // Your trigger channel ID

	// User joins the trigger channel
	if (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID != triggerChannelID) && v.ChannelID == triggerChannelID {
		time.Sleep(1 * time.Second) // Delay to ensure any previous channel deletion is processed

		user, err := s.User(v.UserID)
		if err != nil {
			fmt.Printf("Error retrieving user: %v\n", err)
			return
		}

		// Create a new voice channel with a unique name
		channelName := fmt.Sprintf("%s's Channel %d", user.Username, time.Now().Unix())
		channel, err := s.GuildChannelCreate(v.GuildID, channelName, discordgo.ChannelTypeGuildVoice)
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
	} else if v.BeforeUpdate != nil && v.ChannelID != triggerChannelID {
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
					fmt.Printf("Error deleting channel: %v\n", err)
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
		fmt.Println(m.Emoji.Name)
		fmt.Println(m.ChannelID)
		fmt.Println(info.TextChannelID)
		if info.TextChannelID == m.ChannelID && m.Emoji.Name == "alert" { // Use the emoji you want to trigger the rename
			fmt.Println("In the if statement")
			// Retrieve the message to get its content for the new channel name
			msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
			if err != nil {
				fmt.Printf("Error retrieving message: %v\n", err)
				return
			}

			fmt.Printf("Message content: %s\n", msg.Content)

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
			} else {
				fmt.Printf("Channel %s renamed to %s based on message content\n", voiceChannelID, newName)
			}
			break // Exit the loop after renaming
		}
	}
}
