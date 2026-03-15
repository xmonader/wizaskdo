package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const (
	defaultModel = "qwen/qwen-2.5-coder-32b-instruct"
	apiURL       = "https://openrouter.ai/api/v1/chat/completions"
	maxTokens    = 500
	temperature  = 0.3
)

type SystemPrompt struct {
	Prompt string
}

var systemTemplate = `You are a terminal assistant. The user will ask you a question about terminal commands, shell scripting, or system tasks.

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
		fmt.Fprintln(os.Stderr, "Usage: wizask <your question>")
		fmt.Fprintln(os.Stderr, "Example: wizask 'how to find large files in current directory'")
		os.Exit(1)
	}

	prompt := strings.Join(os.Args[1:], " ")
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: OPENROUTER_API_KEY environment variable not set")
		fmt.Fprintln(os.Stderr, "Get a key at: https://openrouter.ai/keys")
		os.Exit(1)
	}

	result, err := ask(prompt, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func ask(prompt, apiKey string) (string, error) {
	sysTpl := template.Must(template.New("system").Parse(systemTemplate))
	var sysBuf bytes.Buffer
	if err := sysTpl.Execute(&sysBuf, nil); err != nil {
		return "", fmt.Errorf("template error: %w", err)
	}

	req := Request{
		Model: defaultModel,
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
	httpReq.Header.Set("X-Title", "wizask")

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
