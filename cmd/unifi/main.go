package main

import (
	"fmt"

	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"github.com/kashalls/minecraft-router-sidehook/internal/server"
	"github.com/kashalls/minecraft-router-sidehook/internal/unifi"
)

const banner = `
minecraft-router unifi webhook
version: %s (%s)

`

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	fmt.Printf(banner, Version, Gitsha)
	log.Init()

	config := unifi.InitConfig()
	if config.ApiKey == "" || (config.User == "" || config.Password == "") {
		fmt.Println("No API key or user/password provided. Exiting.")
		return
	}

	client, err := unifi.NewClient(&config)
	if err != nil {
		fmt.Println("Error creating Unifi client:", err)
		return
	}

	main, health := server.Start(unifi.InitServer(client))
	server.ShutdownGracefully(main, health)

	fmt.Println("Neither a server host nor a port was provided. Exiting.")
}
