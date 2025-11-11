# blog-aggregator

A command line tool for aggregating RSS feeds and viewing the posts.

## Installation
Install the latest [Go toolchain](https://golang.org/dl/) and a local Postgres database. 
Install `gator` with:

```bash
go install github.com/KasjanK/blog-aggregator/cmd/gator@latest
```

## Config
Create a `.gatorconfig.json` file in your home directory:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the url with your database connection string.

## Usage
Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

Other commands:

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed that already exists in the database
