package main

import (
	"fmt"

	"github.com/kashalls/minecraft-router-sidehook/internal/discord"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
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
	log.Init()

	config := discord.InitConfig()
	if config.Webhook == "" {
		fmt.Println("No webhook URL provided. Exiting.")
		return
	}

	main, health := server.Start(discord.InitServer(config))
	server.ShutdownGracefully(main, health)

	fmt.Println("Neither a webhook URL nor a bot token was provided. Exiting.")
}
