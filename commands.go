package main

import "errors"

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
		return command(s, cmd)
	}
	return errors.New("Unknown command")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFunctions[name] = f
}
