package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s\n", cmd.Name)
	}

	if err := s.db.ResetDb(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	fmt.Println("Database reset successfully!")

	return nil
}
