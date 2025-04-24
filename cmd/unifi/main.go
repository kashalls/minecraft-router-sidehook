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
	client, err := unifi.NewClient(&config)
	if err != nil {
		fmt.Println("Error creating Unifi client:", err)
		return
	}

	if config.VerifyObjects {
		groups, err := client.FetchNetworkObjects()
		if err != nil {
			fmt.Println("Error fetching network objects:", err)
			return
		}

		ipv4Object := unifi.Find(groups, func(g unifi.NetworkGroup) bool {
			return g.Name == config.IPv4ObjectName
		})

		if ipv4Object == nil {
			fmt.Printf("IPv4 object '%s' not found\n", config.IPv4ObjectName)

			object := unifi.NetworkGroup{
				Name:         config.IPv4ObjectName,
				GroupType:    unifi.IPv4GroupType,
				GroupMembers: []string{config.IPv4DefaultObjectValue},
			}
			if err := client.CreateNetworkObject(object); err != nil {
				fmt.Printf("Error creating IPv4 object '%s': %v\n", config.IPv4ObjectName, err)
			} else {
				fmt.Printf("IPv4 object '%s' created successfully\n", config.IPv4ObjectName)
			}
		}

		ipv6Object := unifi.Find(groups, func(g unifi.NetworkGroup) bool {
			return g.Name == config.IPv6ObjectName
		})

		if ipv6Object == nil {
			fmt.Printf("IPv6 object '%s' not found\n", config.IPv6ObjectName)

			object := unifi.NetworkGroup{
				Name:         config.IPv6ObjectName,
				GroupType:    unifi.IPv6GroupType,
				GroupMembers: []string{config.IPv6DefaultObjectValue},
			}
			if err := client.CreateNetworkObject(object); err != nil {
				fmt.Printf("Error creating IPv6 object '%s': %v\n", config.IPv6ObjectName, err)
			} else {
				fmt.Printf("IPv6 object '%s' created successfully\n", config.IPv6ObjectName)
			}
		}
	}

	main, health := server.Start(unifi.InitServer(client))
	server.ShutdownGracefully(main, health)

	fmt.Println("Neither a server host nor a port was provided. Exiting.")
}
