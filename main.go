package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/KasjanK/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	Name 		string
	Arguments 	[]string
}

type commands struct {
	handlerFunctions	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.handlerFunctions[cmd.Name]
	if ok {
		command(s, cmd)
	} else {
		return errors.New("Unknown command")
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFunctions[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Arguments[0]
	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("username has been set to %s", name)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	
	newState := &state{cfg: &cfg}
	cmds := commands{handlerFunctions: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	userArgs := os.Args

	if len(userArgs) < 2 {
		log.Fatal("invalid number of args")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	
	cmd := command{Name: cmdName, Arguments: cmdArgs}
	cmds.run(newState, cmd)
	
}

