package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	botCreatedChannels = make(map[string]bool)
	mutex              sync.Mutex // Mutex to protect access to botCreatedChannels
)

func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	triggerChannelID := "1211873008806010930" // Replace with your trigger channel ID

	fmt.Println("Current botCreatedChannels map:", botCreatedChannels)

	// Check if a user joins the trigger channel from any other state (including no channel)
	if (v.BeforeUpdate == nil || v.BeforeUpdate.ChannelID != triggerChannelID) && v.ChannelID == triggerChannelID {
		time.Sleep(1 * time.Second) // Short delay to ensure previous channel deletion is processed

		user, err := s.User(v.UserID)
		if err != nil {
			fmt.Printf("Error retrieving user: %v\n", err)
			return
		}

		// Create a new voice channel with a unique name using the user's username and a timestamp
		channelName := fmt.Sprintf("%s Channel %d", user.Username, time.Now().Unix())
		fmt.Printf("Creating new channel: %s\n", channelName)
		channel, err := s.GuildChannelCreate(v.GuildID, channelName, discordgo.ChannelTypeGuildVoice)
		if err != nil {
			fmt.Printf("Error creating voice channel: %v\n", err)
			return
		}
		fmt.Printf("Channel created successfully: %s\n", channel.ID)

		mutex.Lock()
		botCreatedChannels[channel.ID] = true
		mutex.Unlock()

		// Move the user to the newly created voice channel
		err = s.GuildMemberMove(v.GuildID, v.UserID, &channel.ID)
		if err != nil {
			fmt.Printf("Error moving user to new channel: %v\n", err)
			// Delete the channel if the user couldn't be moved
			s.ChannelDelete(channel.ID)
			return
		}
	} else if v.BeforeUpdate != nil && v.ChannelID != triggerChannelID {
		mutex.Lock()
		shouldDelete := botCreatedChannels[v.BeforeUpdate.ChannelID]
		mutex.Unlock()

		if shouldDelete {
			guild, err := s.State.Guild(v.GuildID)
			if err != nil {
				fmt.Printf("Error retrieving guild: %v\n", err)
				return
			}

			isEmpty := true
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == v.BeforeUpdate.ChannelID {
					isEmpty = false
					break
				}
			}

			if isEmpty {
				fmt.Printf("Deleting channel: %s\n", v.BeforeUpdate.ChannelID)
				_, err := s.ChannelDelete(v.BeforeUpdate.ChannelID)
				if err != nil {
					fmt.Printf("Error deleting channel: %v\n", err)
					return
				}
				fmt.Printf("Channel deleted successfully: %s\n", v.BeforeUpdate.ChannelID)

				mutex.Lock()
				delete(botCreatedChannels, v.BeforeUpdate.ChannelID)
				mutex.Unlock()
			}
		}
	}
}
