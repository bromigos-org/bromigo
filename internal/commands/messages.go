package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func BotMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the bot was mentioned as the first part of the message
	if len(m.Mentions) > 0 && m.Mentions[0].ID == s.State.User.ID {
		prefix := strings.TrimSpace(strings.TrimPrefix(m.Content, "<@"+s.State.User.ID+">"))
		switch prefix {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		case "help":
			Help(s, m)
		}
	}
}
