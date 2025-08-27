package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/jj-attaq/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("No such user in database\n")
	} else if err != nil {
		return fmt.Errorf("ERROR: %w\n", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User set to: %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}

	name := cmd.Args[0]

	// sql stuff
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s\n", cmd.Name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not list users: %w", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
