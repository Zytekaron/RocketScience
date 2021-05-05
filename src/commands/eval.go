package commands

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// https://emkc.org/
type EvalRequest struct {
	Lang string `json:"language"`
	Code string `json:"source"`
}

// https://emkc.org/
type EvalResponse struct {
	Ran     bool   `json:"ran"`
	Lang    string `json:"language"`
	Version string `json:"version"`
	Output  string `json:"output"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

var (
	// Eval settings/info
	evalEnabled     = true
	evalCodes       = "ctx"
	evalLanguages   = "js, go, java, kt, rust, c, cpp"
	evalMaxNewLines = 20
	evalMaxLength   = 1900

	// File evaluation templates
	evalTemplates = map[string]string{}
)

func init() {
	evalRegister("c")
	evalRegister("cpp", "cc", "c++")
	evalRegister("go", "golang")
	evalRegister("java")
	evalRegister("js", "javascript", "node", "nodejs")
	evalRegister("kt", "kotlin")
	evalRegister("rs", "rust")
}

func evalRegister(name string, aliases ...string) {
	data, err := ioutil.ReadFile("./templates/eval/" + name)
	if err != nil {
		log.Fatal(err)
	}

	evalTemplates[name] = string(data)
	for _, alias := range aliases {
		evalTemplates[alias] = string(data)
	}
}

func EvalCommand(s *discordgo.Session, m *discordgo.Message, args []string) {
	if !evalEnabled && !developer(m.Author) {
		reply(s, m, "Eval is disabled.")
		return
	}

	if len(args) < 1 {
		reply(s, m, "Usage: `eval <lang> [-r] <code>`\nSupported languages: "+evalLanguages)
		return
	}
	if len(args) < 2 {
		reply(s, m, "Include the code to execute.")
		return
	}

	lang := args[0]
	code := strings.Join(args[1:], " ")

	switch strings.ToLower(lang) {
	case "ctx":
		template, ok := evalTemplates[code]
		if !ok {
			reply(s, m, "Unsupported language for context.\nSupported languages: "+evalLanguages)
		} else {
			reply(s, m, "```"+code+"\n"+template+"\n```")
		}
	case "c":
		doLang(s, m, "c", evalTemplates["c"], code)
	case "cpp", "cc", "c++":
		doLang(s, m, "cpp", evalTemplates["cpp"], code)
	case "go", "golang":
		doLang(s, m, "go", evalTemplates["go"], code)
	case "java":
		doLang(s, m, "java", evalTemplates["java"], code)
	case "js", "javascript", "node", "nodejs":
		code = strings.ReplaceAll(code, "`", "\\`")
		doLang(s, m, "js", evalTemplates["js"], code)
	case "kt", "kotlin":
		doLang(s, m, "kt", evalTemplates["kt"], code)
	case "rs", "rust":
		doLang(s, m, "rs", evalTemplates["rs"], code)
	default:
		reply(s, m, "Unsupported language or code.\nSupported languages: "+evalLanguages+"\nSupported codes: "+evalCodes)
	}
}

func evalRequest(lang, code string) (*EvalResponse, error) {
	data, err := json.Marshal(&EvalRequest{
		Code: code,
		Lang: lang,
	})
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(data)
	res, err := http.Post("https://emkc.org/api/v1/piston/execute", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out *EvalResponse
	return out, json.NewDecoder(res.Body).Decode(&out)
}

func doLang(s *discordgo.Session, m *discordgo.Message, lang, template, code string) {
	if strings.HasPrefix(code, "-") && len(code) > 3 && (code[1] == 'b' || code[1] == 'r') {
		code = code[3:]
	} else {
		code = strings.Replace(template, "{{code}}", code, 1)
	}

	res, err := evalRequest(lang, code)
	if err != nil {
		reply(s, m, "An error occurred when making the request:", err)
		return
	}

	if strings.Count(res.Output, "\n") > evalMaxNewLines || len(res.Output) > evalMaxLength {
		evalHaste(s, m, res, lang)
		return
	}

	reply(s, m, "output: (v"+res.Version+") ```xl\n"+res.Output+"\n```")
}

func evalHaste(s *discordgo.Session, m *discordgo.Message, res *EvalResponse, lang string) {
	url, err := haste(res.Output)
	if err != nil {
		reply(s, m, "Could not publish to hst.sh: "+err.Error())
		return
	}

	reply(s, m, "output: (v"+res.Version+") "+url+"."+lang)
}
