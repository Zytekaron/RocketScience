package events

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandHandler(s, m.Message)
}

func MessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	commandHandler(s, m.Message)
}

func commandHandler(s *discordgo.Session, m *discordgo.Message) {
	args := strings.Split(m.Content, " ")
	cmdName := strings.ToLower(args[0])
	args = args[1:]

	cmd, ok := cmds[cmdName]
	if ok {
		cmd(s, m, args)
	}
}
