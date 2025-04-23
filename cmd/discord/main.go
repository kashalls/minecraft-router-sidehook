package main

import (
	"fmt"

	"github.com/kashalls/minecraft-router-sidehook/cmd/discord/configuration"
	"github.com/kashalls/minecraft-router-sidehook/internal/discord/webhook"
	"github.com/kashalls/minecraft-router-sidehook/internal/server"
)

const banner = `
minecraft-router discord webhook
version: %s (%s)

`

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	fmt.Printf(banner, Version, Gitsha)

	config := configuration.Init()

	if config.Webhook != "" {
		main, health := server.Start(webhook.InitServer())
		server.ShutdownGracefully(main, health)
	}

	if config.Token != "" {
		// Todo: Implement Discord bot functionality
	}

	fmt.Println("Neither a webhook URL nor a bot token was provided. Exiting.")
}
