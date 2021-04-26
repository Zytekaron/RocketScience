package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"strings"
)

var commandReplies = map[string]string{}

// Reply in a channel where a message was sent.
func reply(s *discordgo.Session, m *discordgo.Message, data ...interface{}) (*discordgo.Message, error) {
	var out []string
	for _, e := range data {
		out = append(out, fmt.Sprintf("%v", e))
	}
	content := strings.Join(out, " ")

	msgId, ok := commandReplies[m.ID]
	if !ok {
		msg, err := s.ChannelMessageSend(m.ChannelID, content)
		if err != nil {
			return nil, err
		}
		commandReplies[m.ID] = msg.ID
		return msg, err
	}
	return s.ChannelMessageEdit(m.ChannelID, msgId, content)
}

func developer(s *discordgo.Session, m *discordgo.Message) bool {
	if m.Author.ID == "272659147974115328" {
		return true
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command")
	if err != nil {
		log.Println(err)
	}
	return false
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
