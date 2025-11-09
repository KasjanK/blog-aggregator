package main

import (
	"context"
	"fmt"

	"github.com/KasjanK/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not get user: %w", err)
		}

		return handler(s, cmd, user)
	}
}
