package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/jj-attaq/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

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
		UserID:    user.ID,
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s\n", cmd.Name)
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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	// get feed by url
	feedUrl, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	// unfollow query
	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feedUrl.ID,
		UserID: user.ID,
	}); err != nil {
		return fmt.Errorf("couldn't delete feed_follow record: %w\n", err)
	}

	fmt.Printf("Unfollowed '%s' successfully!\n", url)

	return nil
}
