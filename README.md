# wizask

Ask AI for terminal commands directly from your shell.

## Installation

```bash
go build -o wizask ./...
sudo cp wizask /usr/local/bin/
```

## Usage

```bash
wizask "how to find files larger than 100MB"
wizask "compress all log files older than 7 days"
wizask "show disk usage by directory"
```

## Setup

Get an API key from [OpenRouter](https://openrouter.ai/keys) and set:

```bash
export OPENROUTER_API_KEY=your_key_here
```

Add to your `~/.bashrc` or `~/.zshrc` for persistence.

## Configuration

**Model selection** (optional) - defaults to free Llama 3.2:

```bash
export WIZASK_MODEL=meta-llama/llama-3.2-3b-instruct:free  # Free, fast (default)
export WIZASK_MODEL=qwen/qwen3-coder:free          # Free, coding-focused
export WIZASK_MODEL=minimax/minimax-m2.5:free       # Free, large context
```

See all models: https://openrouter.ai/models

## Example Output

```
$ wizask "find all .log files older than 7 days"

```bash
find /var/log -name "*.log" -mtime +7
```
This finds all .log files in /var/log that are older than 7 days. 
Add -delete to remove them (be careful!).
```

## Model

Defaults to `qwen/qwen3-coder:free` - **free**, coding-specialized model.
Override with `WIZASK_MODEL` env var.

## License

MIT
