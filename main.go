package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jj-attaq/gator/internal/config"
	"github.com/jj-attaq/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	// Set initial state and config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	db, err := sql.Open("postgres", programState.cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	programState.db = dbQueries

	// Initiate and register commands
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("register", handlerRegister)
	cmds.register("login", handlerLogin)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerListFeeds)
	// Requires Login
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Get arguemnts
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	commandName := args[1]
	commandArg := args[2:]

	// Use arguments to create commands
	cmd := command{
		Name: commandName,
		Args: commandArg,
	}

	// Run command
	if err := cmds.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
}
