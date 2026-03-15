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

Uses `qwen/qwen-2.5-coder-32b-instruct` by default - a cheap, capable model for coding tasks.

## License

MIT
