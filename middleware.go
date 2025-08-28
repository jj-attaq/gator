package main

import (
	"context"
	"fmt"

	"github.com/jj-attaq/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not logged in: %w\n", err)
		}

		if err := handler(s, cmd, currentUser); err != nil {
			return fmt.Errorf("Error, command %s not executed: %w\n", cmd.Name, err)
		}

		return nil
	}
}
