GO_FLAGS = -ldflags "-X 'github.com/leighmacdonald/discord_log_relay/relay.BuildVersion=`git describe --abbrev=0`'"

all: linux win

linux:
	@GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o build/linux64/discord_log_relay main.go

win:
	@GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o build/win64/discord_log_relay.exe main.go

image:
	@docker build -t leighmacdonald/discord_log_relay:latest .

dist: linux win
	@zip -j discord_log_relay-`git describe --abbrev=0`-win64.zip build/win64/discord_log_relay.exe LICENSE.md
	@zip -j discord_log_relay-`git describe --abbrev=0`-linux64.zip build/linux64/discord_log_relay LICENSE.md
