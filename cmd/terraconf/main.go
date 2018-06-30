package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform/terraform"
	"github.com/jzbruno/terraconf/pkg/terraconf"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: terraconf [stateFile]")
		os.Exit(1)
	}

	stateFile := os.Args[1]

	f, err := os.Open(stateFile)
	if err != nil {
		log.Fatalf("Failed to open state file, err='%s'", err)
	}

	// Ugh, when reading state Terraform displays a message about lineage.
	log.SetOutput(ioutil.Discard)
	state, err := terraform.ReadState(f)
	log.SetOutput(os.Stdout)

	if err != nil {
		log.Fatalf("Failed to read state file, err='%s'", err)
	}

	s, err := terraconf.GetStateConfigString(state)
	if err != nil {
		log.Fatalf("Failed to generate config for state: %s", err)
	}

	fmt.Println(s)
}
