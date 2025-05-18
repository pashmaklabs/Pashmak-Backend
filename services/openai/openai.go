package services_openai

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"pashmak.com/pashmak/bootstrap"
	oa "pashmak.com/pashmak/models/openai"
)

var (
	ErrInvalidAPIKey = errors.New("invalid API key")
	ErrEmptyMessage  = errors.New("message cannot be empty")
)

type OpenAIService struct {
	AppConfig *bootstrap.AppConfig
	APIKey    string
}

func NewOpenAIService(apiKey string) *OpenAIService {
	if apiKey == "" {
		log.Println("Warning: OpenAI API key is empty")
	}
	return &OpenAIService{
		APIKey: apiKey,
	}
}

// CreateChatCompletion creates a new chat completion with the given model.
func (os *OpenAIService) CreateChatCompletion(model string) (*openai.ChatCompletion, error) {
	if os.APIKey == "" {
		log.Println("Error: OpenAI API key is not set")
		return nil, ErrInvalidAPIKey
	}

	client := openai.NewClient(option.WithAPIKey(os.APIKey))
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Say this is a test")},
			Model:    model,
		},
	)
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}
	return resp, nil
}

// NewChatAgent creates a new agent with a constant system prompt.
func (os *OpenAIService) NewChatAgent(model, systemPrompt string) *oa.ChatAgent {
	return oa.NewChatAgent(systemPrompt, model)
}

// SendMessage adds a user message, gets a response from OpenAI, and updates the conversation history.
func (os *OpenAIService) SendMessage(userMessage string, chatAgent *oa.ChatAgent) (string, error) {
	if userMessage == "" {
		return "", ErrEmptyMessage
	}

	if os.APIKey == "" {
		log.Println("Error: OpenAI API key is not set")
		return "", ErrInvalidAPIKey
	}

	chatAgent.AddUserMessage(userMessage)

	client := openai.NewClient(option.WithAPIKey(os.APIKey))
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: chatAgent.GetMessages(),
			Model:    chatAgent.Model,
		},
	)
	if err != nil {
		log.Printf("Error sending message to OpenAI: %v\nAPI Key length: %d\nModel: %s\nMessages: %+v",
			err, len(os.APIKey), chatAgent.Model, chatAgent.GetMessages())
		return "", fmt.Errorf("failed to get response from OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		log.Println("Error: No choices returned from OpenAI")
		return "", errors.New("no response from OpenAI")
	}

	reply := resp.Choices[0].Message.Content
	chatAgent.AddAssistantMessage(reply)

	return reply, nil
}

// GenerateSQL generates a SQL query based on the natural language input
func (os *OpenAIService) GenerateSQL(schema, query string) (string, error) {
	if os.APIKey == "" {
		log.Println("Error: OpenAI API key is not set")
		return "", ErrInvalidAPIKey
	}

	// Create a new SQL chat agent
	sqlAgent := oa.NewSQLChatAgent(schema, "gpt-4.1")
	sqlAgent.AddUserMessage(query)

	client := openai.NewClient(option.WithAPIKey(os.APIKey))
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: sqlAgent.GetMessages(),
			Model:    sqlAgent.Model,
		},
	)
	if err != nil {
		log.Printf("Error generating SQL: %v", err)
		return "", fmt.Errorf("failed to generate SQL: %w", err)
	}

	if len(resp.Choices) == 0 {
		log.Println("Error: No choices returned from OpenAI")
		return "", errors.New("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
