package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jj-attaq/gator/internal/database"
)

var rssTimeLayouts = []string{
	time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano,
	time.RFC822Z,                      // "02 Jan 06 15:04 -0700"
	time.RFC822,                       // "02 Jan 06 15:04 MST"
	"Mon, 02 Jan 2006 15:04:05 GMT",   // some feeds hardcode GMT
	"Mon, 02 Jan 2006 15:04:05 -0700", // same as RFC1123Z but some libs format without day name issues
	"Mon, 02 Jan 2006 15:04:05 MST",
}

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
	// defer ticker.Stop()
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("agg scrape error: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	// Get the next feed to fetch from the DB.
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("No feeds ready to fetch")
		return nil
	} else if err != nil {
		return fmt.Errorf("couldn't get next feeds to fetch: %w", err)
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
		return fmt.Errorf("Could't collect feed %s: %w\n", feed.ID, err)
	}

	// Iterate over the items in the feed and save them to the posts table
	for _, item := range feedData.Channel.Items {
		publishedAt, err := parsePubTime(item.PubDate)
		if err != nil {
			fmt.Printf("warn: bad pubDate %q: %v\n", item.PubDate, err)
			// continue
			publishedAt = time.Now()
		}

		post, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time: publishedAt,
			},
			FeedID: feed.ID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			continue
		} else if err != nil {
			return fmt.Errorf("Couldn't create post: %w", err)
		}

		fmt.Printf("Post created: %s\n", post.Title)
	}

	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(feedData.Channel.Items))

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 && cmd.Args[0] != "" {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("enter valid integer: %w", err)
		}
		limit = n
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <limit>", cmd.Name)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}
	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
	fmt.Printf("--- %s ---\n", post.Title)
	fmt.Printf("    %v\n", post.Description.String)
	fmt.Printf("Link: %s\n", post.Url)
	fmt.Println("=====================================")
}

func parsePubTime(s string) (time.Time, error) {
	var lastErr error
	for _, layout := range rssTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	return time.Time{}, lastErr
}
