package main

import (
	"log"
	"os"

	"github.com/jj-attaq/gator/internal/config"
)

type state struct {
	Config *config.Config
}

func main() {
	// Set initial state and config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	programState := &state{
		Config: &cfg,
	}

	// Initiate and register commands
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	// Get arguemnts
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	commandName := args[1]
	commandArg := args[2:]

	// Use arguments to create commands
	cmd := command{
		Name: commandName,
		Args: commandArg,
	}

	// Run command
	if err := cmds.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
}
