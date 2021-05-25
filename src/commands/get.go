package commands

import (
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var getEnabled = true
var viewHosts = []string{"hst.sh", "pastebin.com"}

func GetCommand(s *discordgo.Session, m *discordgo.Message, args []string) {
	if !getEnabled && !developer(m.Author) {
		reply(s, m, "Get is disabled.")
		return
	}

	u := strings.Join(args, " ")

	for _, host := range viewHosts {
		if strings.Contains(u, host) && !strings.Contains(u, "/raw/") {
			u = strings.Replace(u, host+"/", host+"/raw/", 1)
		}
	}

	if !viewIsUrl(u) {
		reply(s, m, "Invalid URL")
		return
	}

	res, err := http.Get(u)
	if err != nil {
		reply(s, m, "Request failed: "+err.Error())
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		reply(s, m, "Decoding failed: "+err.Error())
		return
	}

	link := ""
	content := string(b)
	if len(content) > 1992 {
		link, err = haste(content)
		if err != nil {
			reply(s, m, "Error pasting to hst.sh: ", err.Error())
			return
		}
		content = content[:128] + "\n..."
	}

	reply(s, m, link+"\n```\n"+content+"\n```")
}

func viewIsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
