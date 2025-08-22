package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("ussage: %s <name>\n", cmd.Name)
	}

	username := cmd.Args[0]

	if err := s.Config.SetUser(username); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User set to: %s\n", username)

	return nil
}
