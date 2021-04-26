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
	evalLanguages   = "js, go, java, kt, rust, c, cpp"
	evalMaxNewLines = 20
	evalMaxLength   = 1536

	// File evaluation templates
	evalCFile    string
	evalCPPFile  string
	evalGoFile   string
	evalJavaFile string
	evalJSFile   string
	evalKTFile   string
	evalRustFile string
)

func init() {
	evalRegister("c", &evalCFile)
	evalRegister("cpp", &evalCPPFile)
	evalRegister("go", &evalGoFile)
	evalRegister("java", &evalJavaFile)
	evalRegister("js", &evalJSFile)
	evalRegister("kt", &evalKTFile)
	evalRegister("rs", &evalRustFile)
}

func evalRegister(name string, variable *string) {
	data, err := ioutil.ReadFile("./templates/eval/" + name)
	if err != nil {
		log.Fatal(err)
	}
	*variable = string(data)
}

func EvalCommand(s *discordgo.Session, m *discordgo.Message, args []string) {
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
	case "c", "clang":
		doLang(s, m, "c", evalCFile, code)
	case "cpp", "c++":
		doLang(s, m, "cpp", evalCPPFile, code)
	case "go", "golang":
		doLang(s, m, "go", evalGoFile, code)
	case "java", "j":
		doLang(s, m, "java", evalJavaFile, code)
	case "js", "javascript", "node", "nodejs":
		code = strings.ReplaceAll(code, "`", "\\`")
		doLang(s, m, "js", evalJSFile, code)
	case "kt", "kotlin", "k":
		doLang(s, m, "kt", evalKTFile, code)
	case "rs", "rust":
		doLang(s, m, "rs", evalRustFile, code)
	default:
		reply(s, m, "Unsupported language. Supported languages: "+evalLanguages)
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
		_, _ = reply(s, m, "An error occurred when making the request:", err)
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
