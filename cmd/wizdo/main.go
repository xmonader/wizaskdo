package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"wizask/pkg/llm"
)

const systemPrompt = `You are a terminal assistant. The user will ask you to perform a terminal command or system task.

Your job is to:
1. Provide ONLY the exact command to run (no explanation, no markdown, no code blocks)
2. The command should be safe to execute directly

Rules:
- Output ONLY the command, nothing else
- No backticks, no markdown, no explanations
- Use standard Unix tools when possible
- Keep it simple and direct

Example:
User: find all log files older than 7 days
Assistant: find /var/log -name "*.log" -mtime +7`

func main() {
	var autoYes bool
	flag.BoolVar(&autoYes, "y", false, "Auto-accept and execute without confirmation")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: wizdo [-y] <command description>")
		fmt.Fprintln(os.Stderr, "Example: wizdo find large files")
		fmt.Fprintln(os.Stderr, "         wizdo -y compress old logs")
		os.Exit(1)
	}

	client, err := llm.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\nGet a key at: https://openrouter.ai/keys\n", err)
		os.Exit(1)
	}

	prompt := strings.Join(flag.Args(), " ")
	cmd, err := client.Ask(systemPrompt, prompt, 500, 0.3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Clean up the command
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimPrefix(cmd, "```")
	cmd = strings.TrimPrefix(cmd, "bash")
	cmd = strings.TrimPrefix(cmd, "sh")
	cmd = strings.TrimSuffix(cmd, "```")
	cmd = strings.TrimSpace(cmd)

	fmt.Printf("🔍 Will execute: %s\n", cmd)

	if !autoYes {
		fmt.Print("\nProceed? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
	} else {
		fmt.Println(" (auto-accepted with -y)")
	}

	fmt.Printf("\n🚀 Executing: %s\n\n", cmd)

	exeCmd := exec.Command("sh", "-c", cmd)
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	exeCmd.Stdin = os.Stdin

	if err := exeCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Command failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Done.")
}
