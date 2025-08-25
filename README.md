# clipstash: Your Lightweight Clipboard Manager

> A simple, fast, command-line clipboard manager that never forgets. `clipstash` runs silently in the background, saving
> everything you copy so you can get it back instantly.

The system clipboard only remembers one thing at a time. `clipstash` gives your clipboard a persistent memory.

-----

## Features

* **Persistent History:** Automatically saves a history of text you copy to a local SQLite database.
* **Simple CLI:** A clean and intuitive command-line interface to `list`, `search`, `get`, and `clear` your history.
* **Lightweight Agent:** Runs as a tiny, efficient background process with minimal resource usage.
* **Custom DB Location:** Use the `CLIPSTASH_DB_PATH` environment variable to store your history anywhere you like.
* **Single Binary:** No dependencies to install. Just download the binary and run it.
* **Cross-Platform:** Designed to work on Linux, macOS, and Windows.

-----

## Installation

You can install `clipstash` by downloading a pre-compiled binary or by building it from the source.

#### From Releases (Recommended)

1. Go to the [Releases Page](https://github.com/mohammadrezajavid-lab/clipstash/releases).
2. Download the latest binary for your operating system.
3. Place it somewhere in your `PATH` (e.g., `/usr/local/bin`).

#### From Source

If you have Go installed, you can build and install `clipstash` with a single command:

```bash
go install github.com/mohammadrezajavid-lab/clipstash@latest
```

-----

## Usage

`clipstash` has two parts: the background agent that saves your history and the CLI you use to access it.

### 1\. Start the Agent

First, you need to start the agent. It's best to run this in the background and have it start automatically when you log
in.

```bash
# Run the agent in the background
clipstash &
```

### 2\. Use the CLI

Once the agent is running, you can use the CLI commands from any terminal.

#### List History

To see the last 15 items you've copied:

```bash
$ clipstash list
4: go mod init clipstash
3: https://github.com/mattn/go-sqlite3
2: A very important API key...
1: SELECT id, content FROM history
```

#### Get an Item

To get an item back onto your clipboard, use its ID.

```bash
# This will copy "A very important API key..." back to your clipboard
clipstash get 2
```

Now you can paste it anywhere\!

#### Search History

To search for a specific item:

```bash
$ clipstash search "sqlite"
3: https://github.com/mattn/go-sqlite3
```

#### Clear History

To permanently delete all saved history:

```bash
clipstash clear
```

-----

## Configuration

By default, `clipstash` stores its database in your user's standard configuration directory. If you want to change this (
for example, to sync your history via Dropbox or another service), you can set the `CLIPSTASH_DB_PATH` environment
variable.

```bash
# Example: Store the database in a Dropbox folder
export CLIPSTASH_DB_PATH="/Users/youruser/Dropbox/clipstash/history.db"

# Now run the agent, and it will use the new path
clipstash &
```

-----

## Development Plan

`clipstash` is actively being developed. Here is the planned roadmap:

- [x] **Phase 1 (Core Functionality):**
    - [x] Robust background agent.
    - [x] SQLite-based storage.
    - [x] `list`, `get`, `search`, and `clear` commands.
- [ ] **Phase 2 (Enhanced Experience):**
    - [ ] Interactive search with a Terminal UI (TUI).
    - [ ] Configuration file for custom settings (e.g., history limit).
    - [ ] Option to ignore certain copied content (e.g., passwords).
- [ ] **Phase 3 (Advanced Features):**
    - [ ] Support for copying and previewing images.
    - [ ] End-to-end encryption for the history database.

-----

## Built With

`clipstash` relies on several amazing open-source projects.

* **[golang.design/x/clipboard](https://pkg.go.dev/golang.design/x/clipboard):** The core of this tool is powered by
  this robust library. Its powerful `Watch` functionality for monitoring the clipboard makes `clipstash` possible.
* **[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite):** A pure Go, CGO-free SQLite implementation that
  enables simple cross-compilation for all major platforms.
* **[Go](https://go.dev):** For providing a powerful and simple language for building fast, cross-platform tools.