package events

import (
	"RocketScience/src/commands"
	"github.com/bwmarrin/discordgo"
)

type Command func(*discordgo.Session, *discordgo.Message, []string)

var cmds = map[string]Command{}

func init() {
	register(commands.EvalCommand, "e", "ev", "eval")
	register(commands.PingCommand, "rsping", "rsonline")
}

func register(run Command, uses ...string) {
	for _, use := range uses {
		cmds[use] = run
	}
}
