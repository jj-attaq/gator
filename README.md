# Gator

Gator is a blog aggregator built using Go, postgreSQL, goose, and sqlc.

### Installation

You will need to have at least Go version 1.23.0, and psql (PostgreSQL) 17.5 installed to run gator.

You can then install Gator itself by running:
    go install github.com/jj-attaq/gator@latest

### Config

The Gator config file is a hidden JSON file located in your home directory: 
    ~/.gatorconfig.json

Inside, you will need to place your local postgres database url into the 'db_url' field as follows: protocol://username:password@host:port/database

The second field is used to determine the current user of the program.

### Commands

#### Register 

    gator register <name>

#### Login

    gator login <name>

#### Agg

    gator agg <time interval between requests>

Ex: 
    gator agg 1m

The program will scrape for feeds every minute.

#### AddFeed

    gator addfeed <name of feed> <url>

#### Browse 

    gator browse <number of posts to list>

Browse will default to 2 posts if no number is specified.
