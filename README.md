# Sen de Yaz

[![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![collage](https://img.shields.io/badge/collage-v0.36.0-6558F5)](https://collage.furkanbaytekin.dev/en/)
[![SQLite](https://img.shields.io/badge/SQLite-embedded-003B57?logo=sqlite&logoColor=white)](https://sqlite.org)

> Write the first sentence; let the community bring the rest.

Sen de Yaz ("you write too") is a platform for stories written together:
one writer opens a story with a title, a theme and an opening line, and the
story grows as the community adds entries, one at a time.

## How it works

- A writer starts a story: a title, a theme description (up to 100
  characters) and an opening text (up to 500).
- Anyone writes the next entry (up to 140 characters), but no one writes
  twice in a row: before you continue, someone else has to.
- The author of the last entry can still edit it until someone else
  contributes; once a new entry lands, the edit closes.

## Getting started

```sh
cp .env.example .env      # fill in the keys: openssl rand -hex 32
go mod tidy
collage dev               # http://localhost:3000, reloading as you edit
```

Build and ship:

```sh
collage build             # -> bin/sen-de-yaz
collage export            # -> dist/, static files
```

`main.go` is the entry point `collage dev` and `collage export` use to run
this program; if you rewrite it, keep both working.

## Layout

```
main.go         sets the app up: plugins, cache, static files, commands
routes.go       where every page and action is registered
pages/          page definitions: builder chains, layouts, action wiring
fragments/      templates and view layer: layouts/ for skeletons, pages/ for content blocks
actions/        POST handlers: validation, service calls, redirects
data/           domain model and services: users/, stories/, sqlite setup
guards/         pre-request checks (e.g. a required session)
utils/      small shared helpers (photo URLs, env loading)
```

## What it is built on

[Collage](https://collage.furkanbaytekin.dev/en/), a Go web framework that
composes pages from fragments and unifies static, dynamic and incremental
rendering in one small core. Sen de Yaz leans on its
[plugins](https://collage.furkanbaytekin.dev/en/docs/plugins/) for
everything from form validation to CSS minification: `validate`, `flash`,
`honeypot`, `secure`, `ratelimit`, `meta`, `jsonld`, `sitemap`,
`opti-image`, `minimizer` and more. Data lives in SQLite through a pure-Go
driver.
