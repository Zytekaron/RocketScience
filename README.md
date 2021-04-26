# RocketScience

A testing platform to evaluate arbitrary code through a Discord bot.

<b>Version:</b> 1.0.0

## Setup
- Install Go from https://golang.org/dl (tested: v1.16.3)
- `git clone https://github.com/Zytekaron/RocketScience`
- `cd RocketScience`
- `go build -o RocketScience main.go` (RocketScience.exe on Windows)
- Create a Discord bot application from the developer portal,
  make it a bot, and copy the token
- create a `.env` file and insert `ROCKET_SCIENCE_TOKEN=BotTokenHere`
- run `./RocketScience` (Unix) or `.\RocketScience.exe` (Windows)

## License
<b>RocketScience</b> is licenced under the [MIT License](https://github.com/Zytekaron/RocketScience/blob/master/LICENSE)
