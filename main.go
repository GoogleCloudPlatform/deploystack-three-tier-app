package main

import (
	"log"

	"github.com/GoogleCloudPlatform/deploystack"
)

func main() {
	deploystack.ClearScreen()
	f := deploystack.HandleFlags()
	s := deploystack.NewStack()
	s.ProcessFlags(f)

	if err := s.ReadConfig("deploystack.json", "deploystack.txt"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if err := s.Process("terraform.tfvars"); err != nil {
		log.Fatalf("problemn collecting the configurations: %s", err)
	}
}
