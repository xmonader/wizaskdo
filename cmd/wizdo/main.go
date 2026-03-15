package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	defaultModel = "liquid/lfm2-8b-a1b"
	apiURL       = "https://openrouter.ai/api/v1/chat/completions"
	maxTokens    = 500
	temperature  = 0.3
)

var systemTemplate = `You are a terminal assistant. The user will ask you to perform a terminal command or system task.

Your job is to:
1. Provide ONLY the exact command to run (no explanation, no markdown, no code blocks)
2. The command should be safe to execute directly
3. If the request is dangerous (rm, chmod, kill, etc.), still provide the command but make it as safe as possible

Rules:
- Output ONLY the command, nothing else
- No backticks, no markdown, no explanations
- Use standard Unix tools when possible
- Keep it simple and direct

Example:
User: find all log files older than 7 days
Assistant: find /var/log -name "*.log" -mtime +7`

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: wizdo <command description>")
		fmt.Fprintln(os.Stderr, "Example: wizdo find large files")
		fmt.Fprintln(os.Stderr, "         wizdo compress old logs")
		os.Exit(1)
	}

	prompt := strings.Join(os.Args[1:], " ")
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: OPENROUTER_API_KEY environment variable not set")
		fmt.Fprintln(os.Stderr, "Get a key at: https://openrouter.ai/keys")
		os.Exit(1)
	}

	model := os.Getenv("WIZASK_MODEL")
	if model == "" {
		model = defaultModel
	}

	cmd, err := getCommand(prompt, apiKey, model)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Clean up the command (remove any markdown artifacts)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimPrefix(cmd, "```")
	cmd = strings.TrimPrefix(cmd, "bash")
	cmd = strings.TrimPrefix(cmd, "sh")
	cmd = strings.TrimSuffix(cmd, "```")
	cmd = strings.TrimSpace(cmd)

	fmt.Printf("🔍 Will execute: %s\n\n", cmd)
	fmt.Print("Proceed? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}

	fmt.Printf("\n🚀 Executing: %s\n\n", cmd)

	// Execute the command
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

func getCommand(prompt, apiKey, model string) (string, error) {
	sysTpl := template.Must(template.New("system").Parse(systemTemplate))
	var sysBuf bytes.Buffer
	if err := sysTpl.Execute(&sysBuf, nil); err != nil {
		return "", fmt.Errorf("template error: %w", err)
	}

	req := Request{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: sysBuf.String()},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/wizask")
	httpReq.Header.Set("X-Title", "wizdo")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return result.Choices[0].Message.Content, nil
}
