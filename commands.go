package main

import (
	"fmt"
	"log"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	runCmd, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("ERROR: '%s' is not a registered command\n", cmd.Name)
	}

	if err := runCmd(s, cmd); err != nil {
		// exit code 1 on command error
		log.Fatal(err)
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
