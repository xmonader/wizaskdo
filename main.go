package main

import (
	"fmt"
	"os"
	"strings"

	"wizask/pkg/llm"
)

const systemPrompt = `You are a terminal assistant. The user will ask you a question about terminal commands, shell scripting, or system tasks.

Your job is to:
1. Provide the exact command to run (in a code block)
2. Explain what the command does concisely

Rules:
- Prefer simple, safe commands
- Use standard Unix tools when possible
- If the command could be destructive, warn the user
- Keep explanations brief and practical
- Format: first show the command in a code block, then explain

Example response:
` + "```bash" + `
find /var/log -name "*.log" -mtime +7 -delete
` + "```" + `
This finds all .log files in /var/log older than 7 days and deletes them. Be careful - this is destructive.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: wizask <your question>")
		fmt.Fprintln(os.Stderr, "Example: wizask find large files")
		os.Exit(1)
	}

	client, err := llm.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\nGet a key at: https://openrouter.ai/keys\n", err)
		os.Exit(1)
	}

	prompt := strings.Join(os.Args[1:], " ")
	result, err := client.Ask(systemPrompt, prompt, 500, 0.3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
