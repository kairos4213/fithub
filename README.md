# FitHub: A web app that assists in tracking and logging fitness metrics

Fitness tracker with big dreams of making health and wellness a little more accessible to everyone.

FitHub is in active development 🏗 -- and will most likely change over time.

## Motivation

Having previously studied and worked in the fitness industry, the biggest hurdle for most people getting started,
was that the cost and lack of knowledge was too great to overcome.

My hope with this application, is to gradually make it a little easier for anyone to learn about and improve their fitness,
but more importantly their health.

## Quick Start

Simply go to the [website](https://fithub.fly.dev) and start exploring!

## Usage

### What It (Currently) Does

* Log body weight, body fat percentage, & muscle mass 💪
* Log personal fitness / health goals 📊
* Create and schedule workouts 🎽
* Search and select exercises from database 🏋‍♂
* Log actual reps, sets, and weights completed for exercises in workout 🥵
* Workout templates to follow along with, or be a starting place for custom
workouts!📘

## Contributing

### Prerequisites

* Go 1.25+
* Node.js / npm (for Tailwind CSS CLI via npx)
* PostgreSQL running locally
* Goose (`go install github.com/pressly/goose/v3/cmd/goose@latest`) for DB migrations
* sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`) for generating database code

> **Note:** `templ` and `air` are managed via `go tool` in go.mod — no manual install needed. They're invoked automatically by `make dev`.

### Clone the repo and install dependencies

```bash
git clone git@github.com:kairos4213/fithub.git && cd fithub
go mod download
```

### Create .env file

```.env
DATABASE_URL=postgres://<user>:<password>@localhost:5432/fithub?sslmode=disable
PORT=8080
FILEPATH_ROOT=./static
TOKEN_SECRET=<generate-a-base64-secret>
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://<user>:<password>@localhost:5432/fithub
GOOSE_MIGRATION_DIR=./sql/schema/

# Optional — Google OAuth (app works without these)
BASE_URL=http://localhost:
GOOGLE_CLIENT_ID=<oauth-client-id>
GOOGLE_CLIENT_SECRET=<oauth-client-secret>
```

### Create the database & run migrations

```bash
createdb fithub
goose up
```

If there are any issues running `goose up` try manually running

`goose -dir sql/schema postgres "$DATABASE_URL" up`

### Generate the database queries

(Only necessary if you make changes to `sql/queries`)

```bash
sqlc generate
```

### Start the dev server (Runs Tailwind, templ, air for hot-reload, and static sync concurrently)

```bash
make dev
```

In your browser, navigate to `http://127.0.0.1:7331/` proxy server for live reloading.

`http://localhost:8080/` will be where the server is actually running.

### Run Tests

```bash
go test ./...
```

### Submit a pull request

If you'd like to contribute, please fork the repo and open a pull request to the `main` branch!

Check out some of the things being worked on below, or potential ideas for future implementation

Found a bug or have a feature idea? [Open an issue](https://github.com/kairos4213/fithub/issues)

## Items Currently Being Worked On

* Create tutorial for users to better understand functionality 📚
* Integrate users Oauth provider photos 📷
* Embed video instructions for workouts 📺
* Expanded exercise selection 📈
* Expand tracking measures beyond simple logging 🧮

## Potential Ideas for Future Implementation

* AI generated workouts 🤖
* Integrations with third party fitness trackers (Fitbit, Apple Watch, etc) ⌚
