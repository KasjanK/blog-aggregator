package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KasjanK/blog-aggregator/internal/config"
	"github.com/KasjanK/blog-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db 	*database.Queries
	cfg *config.Config
}

type command struct {
	Name 		string
	Arguments 	[]string
}

type commands struct {
	handlerFunctions	map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title 		string   	`xml:"title"`
		Link 		string   	`xml:"link"`
		Description string  	`xml:"description"`
		Item 		[]RSSItem   `xml:"item"`
	} `xml:"channel"`
}	

type RSSItem struct {
	Title 		string `xml:"title"` 
	Link 		string `xml:"link"` 
	Description string `xml:"description"` 
	PubDate 	string `xml:"pubDate"` 
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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}
	fmt.Println("successfully deleted all users")
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

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error getting body: %w", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error unmarshaling body: %w", err)
	}
	
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
	
	return &feed, nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	
	db, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("error connecting to db :v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{db: dbQueries, cfg: &cfg}

	cmds := commands{handlerFunctions: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)	
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	userArgs := os.Args

	if len(userArgs) < 2 {
		log.Fatal("invalid number of args")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	
	cmd := command{Name: cmdName, Arguments: cmdArgs}
	if err := cmds.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
	
}

