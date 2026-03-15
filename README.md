# tuifeed Terminal RSS Reader

tuifeed is a terminal RSS feed reader built with Go. This is a pet project — a personal exercise in building TUI interfaces. The code is open for anyone to browse, learn from, or take inspiration.

## What it does

Fetches your RSS feeds and lets you read them without leaving the terminal. Three levels of navigation:

1. **Feed list** — all your configured RSS channels
2. **Article list** — articles from the selected feed
3. **Article view** — fetches the full article from the web and renders it as Markdown in the terminal

## Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — TUI framework, implements The Elm Architecture
- **[Glamour](https://github.com/charmbracelet/glamour)** — Markdown rendering in the terminal
- **[html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown)** — converts fetched HTML articles to Markdown
- Standard library only for RSS XML parsing (`encoding/xml`) and HTTP (`net/http`)

## Getting started

**Prerequisites:** Go 1.21+, internet access.

```bash
git clone git@github.com:Nazar910/tuifeed.git
cd tuifeed
```

Create the feeds config file:

```bash
mkdir -p data
```

Add your feeds to `data/feeds.json`, e.g. you want to follow Andrew Kelley (creator of Zig) and software architecture topics:

```json
[
    "https://andrewkelley.me/rss.xml",
    "https://microservices.io/feed.xml",
    "https://feed.infoq.com/development/"
]
```

Run:

```bash
go run .
```

Or build a binary:

```bash
go build -o tuifeed .
./tuifeed
```

> The program must be run from the project root — it looks for `data/feeds.json` relative to the working directory.

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `Down` | Move cursor down |
| `k` / `Up` | Move cursor up |
| `Enter` | Select / drill in |
| `Esc` | Go back |
| `q` / `Ctrl+C` | Quit |

## Project structure

```
terminal-rss-reader-go/
├── main.go        # TUI model, update loop, rendering
├── feed.go        # RSS fetching, parsing, article rendering
├── go.mod
└── data/
    └── feeds.json # Your feed URLs (gitignored, local only)
```

`data/` is gitignored so your feed list stays local.

## Notes on the implementation

A few things worth pointing out if you're here to learn:

- **Elm Architecture in Go** — Bubble Tea's model/update/view pattern maps cleanly onto a terminal app. `main.go` is a small example of this.
- **Async I/O with tea.Cmd** — feed fetching and article loading run as async commands wrapped in `func() tea.Msg` closures, which is the idiomatic Bubble Tea way to handle side effects without blocking the event loop.
- **Manual viewport** — instead of using Bubble Tea's built-in `viewport` Bubble, article scrolling is implemented as a simple `start`/`end` integer pair on a struct. Straightforward and easy to follow.
- **RSS without a library** — the project parses RSS 2.0 using only `encoding/xml` and struct tags. Works fine for standard feeds; won't handle Atom or media extensions.
- **HTML → Markdown → terminal** — article pages are fetched as raw HTML, converted to Markdown, then rendered with Glamour's dark theme. The rendering chain is easy to swap out if you want different styling.
