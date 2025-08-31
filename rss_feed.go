package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w\n", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}

	// // when running agg:
	// // Error: Unexpected content type: application/rss+xml; charset=utf-8
	// contentType := resp.Header.Get("Content-Type")
	// if contentType != "application/xml" && contentType != "text/xml" {
	// 	return nil, fmt.Errorf("Unexpected content type: %s\n", contentType)
	// }
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") &&
		!strings.Contains(contentType, "text/xml") &&
		!strings.Contains(contentType, "application/rss+xml") {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}

	var feed RSSFeed
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&feed); err != nil {
		return &feed, fmt.Errorf("Error decoding XML: %w\n", err)
	}

	result := formatFeed(feed)

	return &result, nil
}

func formatFeed(feed RSSFeed) RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	items := feed.Channel.Items
	for i, item := range items {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		items[i] = item
		// items[i].Title = html.UnescapeString(el.Title)
		// items[i].Description = html.UnescapeString(el.Description)
	}

	return feed
}
