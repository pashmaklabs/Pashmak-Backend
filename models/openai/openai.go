package models

import (
	"fmt"

	"github.com/openai/openai-go"
)

const SystemPrompt = `
You are a SQL generator for a Persian-language place suggestion app. Given a Persian user query (as a string) and optional latitude/longitude coordinates, generate a PostgreSQL SQL query using these rules:

1. **Always include the columns**: osm_id, name, amenity, and location (using ST_AsText(way) AS location) in the SELECT clause.
2. **For latitude and longitude**, use:
   - ST_Y(ST_Transform(way, 4326)) as latitude
   - ST_X(ST_Transform(way, 4326)) as longitude
3. Use the planet_osm_point table.
4. If a known place type is mentioned (like cafe, restaurant, pharmacy, etc.), filter with the corresponding amenity value.
5. Extract all Persian keywords (e.g. place name, city, area) from the query string.
6. Normalize and tokenize the Persian keywords.
7. Use full-text search with the 'simple' text search configuration:
   to_tsvector('simple', name) @@ to_tsquery('simple', <normalized_query>)
8. Also apply ILIKE filters for each keyword on both name and place columns:
   - name ILIKE '%<keyword>%'
   - place ILIKE '%<keyword>%'
9. Combine the full-text and ILIKE conditions using OR.
10. **If latitude and longitude are provided**, **use these values to filter points within 5000 meters**:
    ST_DWithin(way, ST_Transform(ST_SetSRID(ST_MakePoint([LONGITUDE], [LATITUDE]), 4326), 3857), 5000)
11. Always limit the results to the top 10 using LIMIT 10.
12. Ensure that the generated query includes **osm_id**, **latitude**, and **longitude** using the ST_Y and ST_X functions.
13. Return only the SQL query as a string—no comments or explanation.
14. Filter spatial data using the provided latitude and longitude coordinates for 50000 meters.
`

// ChatAgent represents an agent with a constant system prompt and conversation history.
type ChatAgent struct {
	SystemPrompt string
	Messages     []openai.ChatCompletionMessageParamUnion
	Model        string
}

// SQLChatAgent represents a specialized agent for generating SQL queries
type SQLChatAgent struct {
	ChatAgent
	TableSchema string // Store the table schema for context
}

// NewChatAgent creates a new ChatAgent with the given system prompt and model.
func NewChatAgent(systemPrompt, model string) *ChatAgent {
	return &ChatAgent{
		SystemPrompt: systemPrompt,
		Messages:     []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(systemPrompt)},
		Model:        model,
	}
}

// NewSQLChatAgent creates a new SQLChatAgent with the given schema and model
func NewSQLChatAgent(schema, model string) *SQLChatAgent {
	return &SQLChatAgent{
		ChatAgent: ChatAgent{
			SystemPrompt: "You are a SQL expert. Generate only SQL queries based on the given schema. Do not include any explanations, just the SQL query.",
			Messages:     []openai.ChatCompletionMessageParamUnion{openai.SystemMessage("You are a SQL expert. Generate only SQL queries based on the given schema. Do not include any explanations, just the SQL query.")},
			Model:        model,
		},
		TableSchema: schema,
	}
}

// AddUserMessage adds a user message to the conversation history.
func (a *ChatAgent) AddUserMessage(message string) {
	a.Messages = append(a.Messages, openai.UserMessage(message))
}

// AddAssistantMessage adds an assistant message to the conversation history.
func (a *ChatAgent) AddAssistantMessage(message string) {
	a.Messages = append(a.Messages, openai.AssistantMessage(message))
}

// ClearMessages clears all messages except the system prompt.
func (a *ChatAgent) ClearMessages() {
	a.Messages = []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(a.SystemPrompt)}
}

// GetMessages returns the current conversation history.
func (a *ChatAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return a.Messages
}

// AddUserMessage adds a user message with schema context
func (a *SQLChatAgent) AddUserMessage(message string) {
	// Include schema context with the user message
	fullMessage := fmt.Sprintf("Schema: %s\nQuery: %s", a.TableSchema, message)
	a.Messages = append(a.Messages, openai.UserMessage(fullMessage))
}
