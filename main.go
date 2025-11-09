package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/KasjanK/blog-aggregator/internal/config"
	"github.com/KasjanK/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db 	*database.Queries
	cfg *config.Config
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
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
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

