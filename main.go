package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type Attachment struct {
	Filename string
	MimeType string
	AddedAt  time.Time
}

// readFileContent reads the content of the file at the given path.
// It returns the content as a string, the MIME type of the file, and an error if any occurred.
// If the file is a binary file, the content is returned as a base64 encoded string.
// If the file is a text file, the content is returned as a plain string.
//
// Parameters:
//   - path: The path to the file to be read.
//
// Returns:
//   - string: The content of the file, either as plain text or base64 encoded string.
//   - string: The MIME type of the file.
//   - error: An error if any occurred during reading the file.
func readFileContent(path string) (string, string, error) {
	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	// Detect MIME type
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Fallback to detection by content
		mimeType = http.DetectContentType(content)
	}

	// Handle binary vs text files
	if !strings.HasPrefix(mimeType, "text/") {
		// Convert binary files to base64
		return base64.StdEncoding.EncodeToString(content), mimeType, nil
	}

	// Return text files as string
	return string(content), mimeType, nil
}

func listAttachments(attachments []Attachment) {
	if len(attachments) == 0 {
		fmt.Println("No attachments in context")
		return
	}
	fmt.Println("Active attachments:")
	for i, a := range attachments {
		fmt.Printf("%d. %s (%s) - added %s\n",
			i+1, a.Filename, a.MimeType,
			a.AddedAt.Format("15:04:05"))
	}
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		return
	}

	client := openai.NewClient(apiKey)
	reader := bufio.NewReader(os.Stdin)
	messages := make([]openai.ChatCompletionMessage, 0)
	attachments := make([]Attachment, 0)

	fmt.Println("Start chatting with ChatGPT")
	fmt.Println("Commands:")
	fmt.Println("  /attach <filepath> - Upload a file for context")
	fmt.Println("  /list - Show active attachments")
	fmt.Println("  /quit - Exit the program")

	// main loop - terminated by /quit command
	for {
		fmt.Print("\nYou: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "/quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Handle file attachment
		if strings.HasPrefix(input, "/attach ") {
			filepath := strings.TrimPrefix(input, "/attach ")
			content, mimeType, err := readFileContent(filepath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				continue
			}

			attachments = append(attachments, Attachment{
				Filename: filepath,
				MimeType: mimeType,
				AddedAt:  time.Now(),
			})

			// Add file content to messages
			fileContext := fmt.Sprintf("File content (%s):\n%s", mimeType, content)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    "user",
				Content: fileContext,
			})
			fmt.Printf("File attached: %s\n", filepath)
			continue
		}

		if input == "/list" {
			listAttachments(attachments)
			continue
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "user",
			Content: input,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT4oLatest,
				Messages: messages,
			},
		)

		if err != nil {
			fmt.Printf("Error getting response: %v\n", err)
			// Print token usage statistics
			fmt.Printf("\n|-> Token usage (in error case) - Prompt: %d, Completion: %d, Total: %d\n",
				resp.Usage.PromptTokens,
				resp.Usage.CompletionTokens,
				resp.Usage.TotalTokens)
			continue
		}

		assistantResponse := resp.Choices[0].Message.Content
		fmt.Printf("\nAssistant: %s\n", assistantResponse)

		// Print token usage statistics
		fmt.Printf("\n|-> Token usage - Prompt: %d, Completion: %d, Total: %d\n",
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
			resp.Usage.TotalTokens)

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: assistantResponse,
		})
	}
}
