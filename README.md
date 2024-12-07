# Simple ChatGPT Bot

This project is a Go-based command-line application that allows users to chat with ChatGPT and attach files for context. The application uses the OpenAI API to interact with ChatGPT and supports attaching text and binary files.

## Features

- Chat with ChatGPT via the command line.
- Attach files to provide additional context to the conversation.
- Supports both text and binary files (binary files are base64 encoded).
- Lists active attachments.
- Displays token usage statistics.

## Requirements

- Go 1.23.3 or later
- OpenAI API key

## Installation

1. Clone the repository:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set the `OPENAI_API_KEY` environment variable:
    ```sh
    export OPENAI_API_KEY=<your-openai-api-key>
    ```

## Usage

1. Run the application:
    ```sh
    go run main.go
    ```

2. Start chatting with ChatGPT. Available commands:
    - `/attach <filepath>`: Upload a file for context.
    - `/list`: Show active attachments.
    - `/quit`: Exit the program.

## File Handling

- The `readFileContent` function reads the content of a file and returns it as a string along with its MIME type. Binary files are base64 encoded.
- The `listAttachments` function lists all active attachments with their filenames, MIME types, and the time they were added.

## Example

```sh
Start chatting with ChatGPT
Commands:
  /attach <filepath> - Upload a file for context
  /list - Show active attachments
  /quit - Exit the program

You: Hello, ChatGPT!
Assistant: Hello! 😊 How can I assist you today?
|-> Token usage - Prompt: 12, Completion: 10, Total: 22

You: /attach file_example.xlsx
File attached: file_example.xlsx 

You: /list
Active attachments:
1. file_example.xlsx (application/vnd.openxmlformats-officedocument.spreadsheetml.sheet) - added 19:22:26

You: /quit
Goodbye!