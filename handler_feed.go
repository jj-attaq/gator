package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/jj-attaq/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("could not create feed %w\n", err)
	}

	newFeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	addedFeed, err := s.db.CreateFeedFollow(context.Background(), newFeedParams)
	if err != nil {
		return fmt.Errorf("could not create new feed_follow record: %w\n", err)
	}

	fmt.Printf("Feed '%s' created successfully:\n", addedFeed.FeedName)
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage %s: follow <url>", cmd.Name)
	}
	url := cmd.Args[0]
	// determine current user
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w\n", err)
	}

	// look up the url in args using a new query
	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't find feed from provided url: %w\n", err)
	}

	// create feed follow record
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't create feed_follows record: %w\n", err)
	}

	fmt.Printf("* User:          %s\n", feedFollow.UserName)
	fmt.Printf("* Feed:          %s\n", feedFollow.FeedName)
	fmt.Println("")

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s\n", cmd.Name)
	}

	name := s.cfg.CurrentUserName

	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("could't locate user: %w\n", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't locate feed_follows of user %s\nError: %w\n", user.Name, err)
	}

	fmt.Printf("User %s follows the following feeds:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("* Name:          %s\n", feed.FeedName)
	}
	fmt.Println("")
	fmt.Println("=====================================")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeedAndCreator(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w\n", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* Creator:       %s\n", feed.UserName)
		fmt.Println("")
	}
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
