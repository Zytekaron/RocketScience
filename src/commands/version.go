package commands

import "github.com/bwmarrin/discordgo"

func VersionCommand(s *discordgo.Session, m *discordgo.Message, _ []string) {
	reply(s, m, "version 1.0")
}
