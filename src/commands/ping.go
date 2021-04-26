package commands

import (
	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, m *discordgo.Message, _ []string) {
	reply(s, m, "pong")
}
