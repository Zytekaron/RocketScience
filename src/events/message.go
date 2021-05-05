package events

import (
	"fmt"
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
	if m.Author != nil && m.Author.Bot {
		return
	}

	args := strings.Split(m.Content, " ")
	cmdName := strings.ToLower(args[0])
	args = args[1:]

	cmd, ok := cmds[cmdName]
	if ok {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					desc := "Panic caught in " + cmdName + ": " + fmt.Sprint(err)
					fmt.Println(desc)
					s.ChannelMessageSend(m.ChannelID, desc)
				}
			}()
			cmd(s, m, args)
		}()
	}
}
