package main

import (
	"log"
	"os"

	"github.com/jj-attaq/gator/internal/config"
)

func main() {
	// Initiate and register commands
	var cmds commands
	cmds.registeredCommands = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)

	// Set initial state and config
	s := state{}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s.Config = &cfg

	// Get arguemnts
	args := os.Args
	if len(args) < 2 {
		log.Fatal("ERROR: command and arguements required")
	}

	commandName := args[1]
	commandArg := args[2:]

	// Use arguments to create commands
	cmd := command{
		name:      commandName,
		arguments: commandArg,
	}

	// Run command
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
