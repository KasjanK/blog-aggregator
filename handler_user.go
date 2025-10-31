package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KasjanK/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Arguments[0]

	getUser, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldnt find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("username has been set to %s", getUser.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Arguments[0] 
	newUser, err := s.db.CreateUser(context.Background(), 
	database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	},
)
if err != nil {
	return fmt.Errorf("could not create user: %w", err)
}

err = s.cfg.SetUser(newUser.Name)
if err != nil {
	return fmt.Errorf("could not set current user: %w", err)
}

fmt.Printf("user was created. ID: %d, created at: %s, updated at: %s, name: %s", newUser.ID, newUser.CreatedAt, newUser.UpdatedAt, newUser.Name)

return nil
}

func handlerUsers(s *state, cmd command) error {
	usersList, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users: %w", err)
	}

	for _, user := range usersList {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}
