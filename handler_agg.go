package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jj-attaq/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("agg scrape: %w", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	// Get the next feed to fetch from the DB.
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Couldn't get next feeds to fetch: %w", err)
	}

	fmt.Println("Found a feed to fetch!")
	if err := scrapeFeed(s, nextFeed); err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	return nil
}

func scrapeFeed(s *state, feed database.Feed) error {
	ctx := context.Background()

	// Mark it as fetched.
	_, err := s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("Couldn't mark feed %s fetched: %w", feed.Name, err)
	}

	// Fetch the feed using the URL (we already wrote this function)
	feedData, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("Could't collect feed %s: %w\n", err)
	}

	// Iterate over the items in the feed and print their titles to the console.
	for _, item := range feedData.Channel.Items {
		fmt.Printf("%s\n", item.Title)
	}

	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(feedData.Channel.Items))
	return nil
}
