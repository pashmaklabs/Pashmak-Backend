package services_place

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models "pashmak.com/pashmak/models/openai"
	oa "pashmak.com/pashmak/models/openai"
	sp "pashmak.com/pashmak/serializers/place"
	services_openai "pashmak.com/pashmak/services/openai"
)

type PlaceService struct {
	DB            *gorm.DB
	AppConfig     *bootstrap.AppConfig
	OpenAIService *services_openai.OpenAIService
}

func NewPlaceService(db *gorm.DB, appconfig *bootstrap.AppConfig, openaiService *services_openai.OpenAIService) *PlaceService {
	return &PlaceService{
		DB:            db,
		AppConfig:     appconfig,
		OpenAIService: openaiService,
	}
}

func (ps *PlaceService) GetPlaceByID(id uint) (*sp.GetPlaceByIDResponse, error) {
	var results []sp.GetPlaceByIDResponse

	query := `
        SELECT 
            osm_id,
            name,
            amenity,
            ST_Y(ST_Transform(way, 4326)) as latitude,
            ST_X(ST_Transform(way, 4326)) as longitude
        FROM planet_osm_point
        WHERE osm_id = ?`

	err := ps.DB.Raw(query, id).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no place found with ID %d", id)
	} else if len(results) > 1 {
		return nil, fmt.Errorf("multiple places found with ID %d", id)
	}

	return &results[0], nil
}

func (ps *PlaceService) SaveSearch(userID uint, sessionID string, loggedIn bool, query string) error {
	// Save to search history
	history := models.SearchHistory{
		UserID:    userID,
		SessionID: sessionID,
		Query:     query,
	}
	if loggedIn {
		history.SessionID = "" // Clear session_id for logged-in users
	}

    // // Use transaction to update or create history
    // err = h.DB.Transaction(func(tx *gorm.DB) error {
    //     var history models.SearchHistory
    //     result := tx.Where("identifier = ?", identifier).First(&history)
    //     if result.Error == gorm.ErrRecordNotFound {
    //         history = models.SearchHistory{
    //             Identifier:  identifier,
    //             IsAnonymous: isAnonymous,
    //             Queries:     []models.QueryEntry{newEntry},
    //         }
    //         return tx.Create(&history).Error
    //     } else if result.Error != nil {
    //         return result.Error
    //     }

    //     history.Queries = append(history.Queries, newEntry)
    //     // Limit to last 50 queries
    //     if len(history.Queries) > 50 {
    //         history.Queries = history.Queries[len(history.Queries)-50:]
    //     }
    //     return tx.Save(&history).Error
    // })
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save search history"})
    //     return
    // }


	if err := ps.DB.Create(history); err != nil{
		return err.Error
	}
	return nil
}

func (ps *PlaceService) SearchPlace(query string) ([]sp.GetPlaceByIDResponse, error) {
	var results []sp.GetPlaceByIDResponse

	// Create a new SQL chat agent with our specialized system prompt
	sqlAgent := oa.NewSQLChatAgent("", "gpt-4.1")
	sqlAgent.SystemPrompt = oa.SystemPrompt // Use our specialized prompt for place search
	sqlAgent.ClearMessages()                // Reset messages to include the new system prompt

	// Add the user's query
	sqlAgent.AddUserMessage(query)

	// Get the generated SQL
	client := openai.NewClient(option.WithAPIKey(ps.AppConfig.OpenaiApiKey))
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: sqlAgent.GetMessages(),
			Model:    sqlAgent.Model,
		},
	)
	if err != nil {
		log.Printf("Error generating SQL: %v\nAPI Key length: %d\nModel: %s\nMessages: %+v",
			err, len(ps.AppConfig.OpenaiApiKey), sqlAgent.Model, sqlAgent.GetMessages())
		// Fallback to simple ILIKE search if SQL generation fails
		fallbackQuery := `
			SELECT 
				name,
				amenity,
				ST_Y(ST_Transform(way, 4326)) as latitude,
				ST_X(ST_Transform(way, 4326)) as longitude
			FROM planet_osm_point
			WHERE name ILIKE ?
			LIMIT 10`

		err = ps.DB.Raw(fallbackQuery, "%"+query+"%").Scan(&results).Error
		if err != nil {
			log.Printf("Fallback search failed: %v", err)
			return nil, fmt.Errorf("fallback search failed: %w", err)
		}
		return results, nil
	}

	if len(resp.Choices) == 0 {
		log.Printf("No SQL query generated from OpenAI response: %+v", resp)
		return nil, fmt.Errorf("no SQL query generated")
	}

	// Execute the generated SQL query
	generatedSQL := resp.Choices[0].Message.Content
	log.Printf("Generated SQL: %s", generatedSQL)

	err = ps.DB.Raw(generatedSQL).Scan(&results).Error
	if err != nil {
		log.Printf("Failed to execute generated SQL: %v\nSQL: %s", err, generatedSQL)
		return nil, fmt.Errorf("failed to execute generated SQL: %w", err)
	}


	return results, nil
}
