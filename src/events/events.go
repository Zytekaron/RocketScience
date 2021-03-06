package events

import (
	"RocketScience/src/commands"
	"github.com/bwmarrin/discordgo"
)

type Command func(*discordgo.Session, *discordgo.Message, []string)

var cmds = map[string]Command{}

func init() {
	register(commands.DevCommand, "rsdev")
	register(commands.EvalCommand, "ev", "eval")
	register(commands.GetCommand, "rsget", "rsview")
	register(commands.HasteCommand, "hst", "haste")
	register(commands.PingCommand, "rsping", "rsonline")
	register(commands.VersionCommand, "rsv", "rsversion")
}

func register(run Command, uses ...string) {
	for _, use := range uses {
		cmds[use] = run
	}
}
