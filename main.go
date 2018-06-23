package main

import (
	"github.com/hashicorp/consul/api"
	"log"
)

func main() {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"

	client, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Init client failed:", err)
	}

	allNodes, _, err := client.Catalog().Nodes(nil)
	if err != nil {
		log.Panicln("Query all known nodes failed:", err)
	}

	allClients := map[string]*api.Client{}
	for _, node := range allNodes {
		tmpConfig := api.DefaultConfig()
		tmpConfig.Address = node.Address + ":8500"
		tmpClient, err := api.NewClient(tmpConfig)
		if err != nil {
			log.Println("Client:", tmpConfig.Address, "create Failed!")
		} else {
			allClients[tmpConfig.Address] = tmpClient
		}
	}

	for address, tmpClient := range allClients {

		allChecks, err := tmpClient.Agent().Checks()
		if err != nil {
			log.Println("Get registered checks failed:", address)
			continue
		}

		log.Println("Clean ===>", address)

		for _, v := range allChecks {
			if v.Status == "critical" {
				log.Println("Deregister ==>", v.ServiceID)
				client.Agent().ServiceDeregister(v.ServiceID)
			}
		}
	}
}
