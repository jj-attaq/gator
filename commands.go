package main

import (
	"fmt"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	runCmd, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("ERROR: '%s' is not a registered command\n", cmd.name)
	}

	if err := runCmd(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	args := cmd.arguments
	if len(args) < 1 {
		return fmt.Errorf("ERROR: No arguments provided to 'login' command\n")
	}
	if len(args) > 1 {
		return fmt.Errorf("ERROR: 'Login' command expects only 1 argument\n")
	}

	username := args[0]

	if err := s.Config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User set to: %s\n", username)

	return nil
}
