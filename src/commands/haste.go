package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func HasteCommand(s *discordgo.Session, m *discordgo.Message, args []string) {
	content := strings.Join(args, " ")
	url, err := haste(content)
	if err != nil {
		reply(s, m, "Could not publish to hst.sh: "+err.Error())
		return
	}

	reply(s, m, url)
}
