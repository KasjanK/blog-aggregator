package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}
	fmt.Println("successfully deleted all users")
	return nil
}
