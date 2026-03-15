# wizask

Ask AI for terminal commands directly from your shell.

[![CI](https://github.com/xmonader/wizask/actions/workflows/ci.yml/badge.svg)](https://github.com/xmonader/wizask/actions/workflows/ci.yml)
[![Release](https://github.com/xmonader/wizask/actions/workflows/release.yml/badge.svg)](https://github.com/xmonader/wizask/actions/workflows/release.yml)

## Tools

### `wizask` - Ask for commands
Returns command + explanation (non-destructive):
```bash
wizask find large files
```

### `wizdo` - Execute commands
Finds command, shows preview, asks confirmation, then executes:
```bash
wizdo compress old log files
```

Auto-accept with `-y`:
```bash
wizdo -y list files
```

## Installation

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/xmonader/wizask/releases):

```bash
# Linux/macOS
curl -L https://github.com/xmonader/wizask/releases/latest/download/wizask_Linux_x86_64.tar.gz | tar xz
sudo mv wizask wizdo /usr/local/bin/

# Or use wget
wget https://github.com/xmonader/wizask/releases/latest/download/wizask_Linux_x86_64.tar.gz
tar xzf wizask_Linux_x86_64.tar.gz
sudo mv wizask wizdo /usr/local/bin/
```

### Homebrew

```bash
brew tap xmonader/tap
brew install wizask
```

### Build from Source

Requires Go 1.21+:

```bash
git clone https://github.com/xmonader/wizask.git
cd wizask
make build
sudo make install
```

Or with `go install`:

```bash
go install github.com/xmonader/wizask@latest
go install github.com/xmonader/wizask/cmd/wizdo@latest
```

## Usage

Quotes are optional for simple queries:

```bash
wizask how to find files larger than 100MB
wizask find large files
wizask "compress all log files older than 7 days"
wizask 'grep "error" in /var/log'
```

Use quotes when your query contains special characters (`*`, `|`, `$`, etc.)

### Flags

**wizask:**
```bash
wizask -v              # Show version
wizask -version        # Show version
```

**wizdo:**
```bash
wizdo -y               # Auto-accept and execute
wizdo -v               # Show version
wizdo -version         # Show version
```

## Setup

Get an API key from [OpenRouter](https://openrouter.ai/keys) and set:

```bash
export OPENROUTER_API_KEY=your_key_here
```

Add to your `~/.bashrc` or `~/.zshrc` for persistence.

## Configuration

**Model selection** (optional) - defaults to ultra-cheap LFM2:

```bash
export WIZASK_MODEL=liquid/lfm2-8b-a1b            # ~$0.00000001/1k (default)
export WIZASK_MODEL=qwen/qwen2.5-coder-7b-instruct # ~$0.00000003/1k, coding
export WIZASK_MODEL=meta-llama/llama-3.2-3b-instruct:free  # Free (rate-limited)
```

See all models: https://openrouter.ai/models

## Example Output

### wizask (ask only)
```
$ wizask find all .log files older than 7 days

```bash
find /var/log -name "*.log" -mtime +7
```
This finds all .log files in /var/log that are older than 7 days.
Add -delete to remove them (be careful!).
```

### wizdo (execute)
```
$ wizdo show current date

🔍 Will execute: date +"%Y-%m-%d"

Proceed? [y/N]: y

🚀 Executing: date +"%Y-%m-%d"
2026-03-15

✅ Done.
```

## Model

Defaults to `liquid/lfm2-8b-a1b` - **~$0.00000001/1k tokens** (extremely cheap).
At this rate, 1000 queries cost less than $0.01.

Override with `WIZASK_MODEL` env var.

## Development

```bash
# Run tests
make test

# Build locally
make build

# Clean
make clean

# Release (requires goreleaser)
make release-snapshot
```

## License

MIT
