module github.com/dangdennis/crossing/bot

go 1.14

require (
	github.com/bwmarrin/discordgo v0.22.0
	github.com/dangdennis/crossing/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.4.1-0.20190629185528-ae1634f6a989 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200824131525-c12d262b63d8 // indirect
	golang.org/x/tools v0.0.0-20200117012304-6edc0a871e69 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/dangdennis/crossing/common => ../common
