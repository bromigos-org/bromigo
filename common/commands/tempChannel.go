package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botCreatedChannels = make(map[string]bool)

func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	triggerChannelID := "1211873008806010930" // Your trigger channel ID

	if v.BeforeUpdate == nil && v.ChannelID == triggerChannelID {
		// Create a new voice channel
		channel, err := s.GuildChannelCreate(v.GuildID, "New Channel", discordgo.ChannelTypeGuildVoice)
		if err != nil {
			fmt.Printf("Error creating voice channel: %v\n", err)
			return
		}
		botCreatedChannels[channel.ID] = true

		// Move the user to the newly created voice channel
		err = s.GuildMemberMove(v.GuildID, v.UserID, &channel.ID)
		if err != nil {
			fmt.Printf("Error moving user to new channel: %v\n", err)
			return
		}
	} else if v.BeforeUpdate != nil {
		// User left a channel, check if it's one of the bot-created channels
		if botCreatedChannels[v.BeforeUpdate.ChannelID] {
			// Delay to handle quick reconnects
			time.Sleep(5 * time.Second)

			// Check if the channel is empty by counting users in the channel
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

			// If the channel is empty, delete it
			if isEmpty {
				ch, err := s.ChannelDelete(v.BeforeUpdate.ChannelID)
				if err != nil {
					fmt.Printf("Error deleting channel: %v\n", err)
					return
				}
				delete(botCreatedChannels, ch.ID)
			}
		}
	}
}
