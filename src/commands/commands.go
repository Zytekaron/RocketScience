package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"strings"
)

var commandReplies = map[string]string{}

// Reply in a channel where a message was sent, and edit the response if called again for the same message
func reply(s *discordgo.Session, m *discordgo.Message, data ...interface{}) (*discordgo.Message, error) {
	msgId, ok := commandReplies[m.ID]
	if !ok {
		msg, err := send(s, m.ChannelID, data...)
		if err != nil {
			return nil, err
		}
		commandReplies[m.ID] = msg.ID
		return msg, err
	}
	return edit(s, m.ChannelID, msgId, data...)
}

func send(s *discordgo.Session, cid string, data ...interface{}) (*discordgo.Message, error) {
	content := join(" ", data...)
	return s.ChannelMessageSend(cid, content)
}

func edit(s *discordgo.Session, cid, mid string, data ...interface{}) (*discordgo.Message, error) {
	content := join(" ", data...)
	return s.ChannelMessageEdit(cid, mid, content)
}

func join(delim string, data ...interface{}) string {
	var out []string
	for _, e := range data {
		out = append(out, fmt.Sprintf("%v", e))
	}
	return strings.Join(out, delim)
}

func developer(u *discordgo.User) bool {
	return u.ID == "272659147974115328"
}

func haste(str string) (string, error) {
	res, err := http.Post("https://hst.sh/documents", "application/json", strings.NewReader(str))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var data struct {
		Key string `json:"key"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return "https://hst.sh/" + data.Key, nil
}
