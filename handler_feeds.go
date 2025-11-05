package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())	
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	for _, feed := range feeds {
		userid, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting userid: %w", err)
		}
		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Created by: %s\n", userid.Name) 
		fmt.Println("=========================================")
	}
	return nil
}
