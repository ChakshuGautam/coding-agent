package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"agent/tools"

	"github.com/joho/godotenv"

	"google.golang.org/genai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system env vars")
	}

	ctx := context.Background()

	agentTools := tools.GetTools()

	apiKey := os.Getenv("GEMINI_API_KEY")

	client, errInClientGenai := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if errInClientGenai != nil {
		log.Fatal(errInClientGenai)
	}

	scanner := bufio.NewScanner(os.Stdin)

	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		message := scanner.Text()
		return message, true
	}

	agent := NewAgent(client, getUserMessage, agentTools)
	err := agent.Run(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File execution ended")
}

func NewAgent(client *genai.Client, getUserMessage func() (string, bool), tools []tools.ToolDefinition) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

type Agent struct {
	client         *genai.Client
	getUserMessage func() (string, bool)
	tools          []tools.ToolDefinition
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []*genai.Part{}
	fmt.Println("Chat with Gemini (use 'ctrl-c' to exit)")

	readUserInput := true

	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			message, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := &genai.Part{Text: message}
			conversation = append(conversation, userMessage)
		}

		response, err := a.runInference(ctx, conversation)
		if err != nil {
			log.Fatal(err)
		}

		agentMessage := &genai.Part{Text: response.Text()}

		toolResults := []*genai.Part{}

		if response.FunctionCalls() != nil {
			functionCalls := response.FunctionCalls()
			if len(functionCalls) > 0 {
				// Process function calls and add FunctionResponse parts to toolResults
				for _, functionCall := range functionCalls {
					var tool *tools.ToolDefinition
					for _, t := range a.tools {
						if t.Definition.Name == functionCall.Name {
							tool = &t
							break
						}
					}
					if tool != nil {
						fmt.Printf("\u001b[92mtool:\u001b[0m %s\n", tool.Definition.Name)
						toolResult, err := tool.Function(functionCall) // function call
						if err != nil {
							log.Fatal(err)
						}
						if len(toolResult) == 0 {
							readUserInput = true
							continue
						}

						// Construct the FunctionResponse part
						toolResponsePart := &genai.Part{
							FunctionResponse: &genai.FunctionResponse{
								Name: functionCall.Name,
								Response: map[string]any{
									"content": toolResult,
								},
							},
						}
						toolResults = append(toolResults, toolResponsePart)
					}
				}
			} else {
				// No function calls, just a text response. Add it to conversation.
				conversation = append(conversation, agentMessage)
				fmt.Printf("\u001b[93mGemini\u001b[0m: %s\n", response.Text())
			}
		} else {
			// No function calls, just a text response. Add it to conversation.
			conversation = append(conversation, agentMessage)
			fmt.Printf("\u001b[93mGemini\u001b[0m: %s\n", response.Text())
		}

		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}

		readUserInput = false
		// Add the results of the function calls to the conversation
		conversation = append(conversation, toolResults...)
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []*genai.Part) (*genai.GenerateContentResponse, error) {
	functionDeclarations := make([]*genai.FunctionDeclaration, len(a.tools))
	for i, tool := range a.tools {
		functionDeclarations[i] = tool.Definition
	}
	response, err := a.client.Models.GenerateContent(ctx, "gemini-2.5-flash-preview-04-17", []*genai.Content{{Parts: conversation}}, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You are a helpful assistant that can use tools to help the user. Always see if you can use a tool to help the user. If you can't, just answer the question."}}},
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: functionDeclarations,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
