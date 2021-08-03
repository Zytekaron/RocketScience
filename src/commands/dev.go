package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

func DevCommand(s *discordgo.Session, m *discordgo.Message, args []string) {
	if !developer(m.Author) {
		reply(s, m, "You are not a developer.")
		return
	}

	if len(args) == 0 {
		reply(s, m, "Invalid")
		return
	}

	option := strings.ToLower(args[0])
	args = args[1:]

	var str string
	switch option {
	case "enable", "e", "disable", "d", "toggle", "t":
		if len(args) == 0 {
			reply(s, m, "Missing option")
			return
		}
		str = devToggle(args[0], rune(option[0]))
	case "stop", "s":
		send(s, m.ChannelID, "RocketScience is shutting down")
		time.Sleep(time.Second) // give it some time to go through
		fallthrough
	case "exitnow", "en", "ei", "kill":
		fmt.Println("Stop command used. Exiting.")
		os.Exit(0)
	default:
		str = "Invalid option: " + option
	}

	reply(s, m, str)
}

func devToggle(command string, option rune) string {
	// yeah, I know this is kinda weird but it's cleaner
	// than all the alternatives I've come up with
	transform := func(in bool) bool {
		if option == 't' {
			return !in
		}
		return option == 'e' // false for 'd'
	}

	switch command {
	case "eval":
		evalEnabled = transform(evalEnabled)
		return "eval is now " + devEnabled(evalEnabled)
	case "get":
		getEnabled = transform(getEnabled)
		return "get is now " + devEnabled(getEnabled)
	}
	return "Unsupported command"
}

func devEnabled(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}
