package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KasjanK/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Arguments[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not get feed %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(), 
		CreatedAt: 	time.Now(), 
		UpdatedAt: 	time.Now(), 
		UserID: 	user.ID, 
		FeedID: 	feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create follow: %w", err)
	}

	fmt.Println("Feed follow created successfully!")
	printFeedFollow(follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	for _, follow := range follows {
		fmt.Println(follow.Name)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

