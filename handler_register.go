package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jj-attaq/gator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}

	name := cmd.Args[0]

	// sql stuff
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	_, err := s.db.GetUser(context.Background(), user.Name)
	if err == nil {
		return fmt.Errorf("User %s is already registered\n", user.Name)
	} else if err == sql.ErrNoRows {
		createdUser, err := s.db.CreateUser(context.Background(), user)
		if err != nil {
			return err
		}

		if err := s.cfg.SetUser(createdUser.Name); err != nil {
			return err
		}

		fmt.Printf("Registered user %s into database\n", createdUser.Name)
	} else {
		return err
	}
	return nil
}
