package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/KasjanK/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Arguments[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit :w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", 	   post.Description.String)
		fmt.Printf("Link: %s\n",   post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
