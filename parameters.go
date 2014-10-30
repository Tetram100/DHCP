package main

import (
	"encoding/json"
	"log"
	"os"
)

type Parameters struct {
	Allocation_time int    `json:"allocationTime"` // In seconds
	Our_Network     string `json:"ourNetwork"`     // In CIDR notation
	IP_server       string `json:"ipServer"`
	IP_DNS          string `json:"ipDns"`
	Netmask         string `json:"netmask"`
	OutPort         int    `json:"out_port"`
}

var parameters Parameters

func getParameters() {

	log.Println("Chargement des paramètres")
	file, err := os.Open("parameters.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&parameters)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Paramètres chargés")
}
