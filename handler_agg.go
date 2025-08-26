package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	// if len(cmd.Args) != 1 {
	// 	return fmt.Errorf("usage: agg <link>\n")
	// }
	// link := cmd.Args[0]
	link := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), link)
	if err != nil {
		return fmt.Errorf("Error: %w\n", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
