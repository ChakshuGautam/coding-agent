# Coding Agent

This is a simple coding agent that uses Gemini to help the user edit code. Adapted from https://ampcode.com/how-to-build-an-agent
The only change I made is to ensure it uses the official [Go client for Gemini](https://github.com/googleapis/go-genai).

## Usage

```bash
go run main.go
```

## Tools

The agent has access to the following tools:

- `editFile`: Edit a file
- `readFile`: Read a file
- `writeFile`: Write to a file
