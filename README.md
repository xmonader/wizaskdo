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

**Model selection** (optional) - defaults to ultra-cheap LFM2:

```bash
export WIZASK_MODEL=liquid/lfm2-8b-a1b            # ~$0.00000001/1k (default)
export WIZASK_MODEL=qwen/qwen2.5-coder-7b-instruct # ~$0.00000003/1k, coding
export WIZASK_MODEL=meta-llama/llama-3.2-3b-instruct:free  # Free (rate-limited)
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

Defaults to `liquid/lfm2-8b-a1b` - **~$0.00000001/1k tokens** (extremely cheap).
Override with `WIZASK_MODEL` env var.

## License

MIT
