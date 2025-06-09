package services_place

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"
	webp "github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	oa "pashmak.com/pashmak/models/openai"
	models_place "pashmak.com/pashmak/models/place"
	sp "pashmak.com/pashmak/serializers/place"
	services_openai "pashmak.com/pashmak/services/openai"
)

var (
	ErrInvalidFile      = errors.New("invalid file type or size")
	ErrNotFound         = errors.New("image not found")
	ErrMinioUnavailable = errors.New("minio unavailable")
	ErrInvalidSize      = errors.New("file too large")
)

type PlaceService struct {
	DB            *gorm.DB
	AppConfig     *bootstrap.AppConfig
	OpenAIService *services_openai.OpenAIService
	Minio         *minio.Client
}

func NewPlaceService(db *gorm.DB, appconfig *bootstrap.AppConfig, openaiService *services_openai.OpenAIService, minioClient *minio.Client) *PlaceService {
	return &PlaceService{
		DB:            db,
		AppConfig:     appconfig,
		OpenAIService: openaiService,
		Minio:         minioClient,
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

func (ps *PlaceService) SaveSearch(userID *uint, sessionID string, loggedIn bool, query string) error {
	history := oa.SearchHistory{
		UserID:    userID,
		SessionID: sessionID,
		Query:     query,
	}
	if loggedIn {
		history.SessionID = "" // Clear session_id for logged-in users
	}

	if err := ps.DB.Create(&history); err != nil{
		return err.Error
	}
	return nil
}


func (ps *PlaceService) SearchPlace(q string, lat string, long string) ([]sp.GetPlaceByIDResponse, error) {
	query := fmt.Sprintf("Query: %s\nLatitude: %s\nLongitude: %s", q, lat, long)
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
			SELECT osm_id, name, NULL AS latitude, NULL AS longitude, NULL as amenity, id
			FROM places
			WHERE name ILIKE ?
			UNION
			SELECT osm_id, name, ST_Y(way) AS latitude, ST_X(way) AS longitude, amenity, NULL AS place_id
			FROM planet_osm_point
			WHERE name ILIKE ? AND osm_id NOT IN (SELECT osm_id FROM places WHERE osm_id IS NOT NULL)
			LIMIT 50`

		err = ps.DB.Raw(fallbackQuery, "%"+q+"%", "%"+q+"%").Scan(&results).Error
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

// validateImage checks file extension and size
func (ps *PlaceService) validateImage(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", ErrInvalidFile
	}
	if file.Size > 1<<24 {
		return "", ErrInvalidSize
	}
	return ext, nil
}

// UploadPlaceImage handles uploading a new image for a place and updating its ImageURLs array.
func (ps *PlaceService) UploadPlaceImage(place models_place.Place, file *multipart.FileHeader) (string, error) {
	_, err := ps.validateImage(file)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	img, _, err := image.Decode(fileReader)
	if err != nil {
		return "", err
	}

	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 30}); err != nil {
		return "", err
	}
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ".webp")
	Reader := bytes.NewReader(buf.Bytes())
	timedCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = ps.Minio.PutObject(
		timedCtx,
		"place-photos",
		objectName,
		Reader,
		Reader.Size(),
		minio.PutObjectOptions{ContentType: "image/webp"},
	)
	if err != nil {
		return "", err
	}

	// Generate public URL
	imgURL := fmt.Sprintf("%s/places/images/%s", ps.AppConfig.ServerHost, objectName)
	if place.Images == nil {
		place.Images = []models_place.Image{}
	}
	place.Images = append(place.Images, models_place.Image{
		PlaceID: place.ID,
		URL: imgURL,
		AltText: "image",
		Caption: "caption",
	})
	if err := ps.DB.Save(&place).Error; err != nil {
		return "", err
	}

	return objectName, nil
}

// GetPlaceImage retrieves an image for a place by image filename
func (ps *PlaceService) GetPlaceImage(placeID uint, imageName string) (io.ReadCloser, string, error) {
	if imageName == "" {
		return nil, "", ErrInvalidFile
	}

	obj, err := ps.Minio.GetObject(context.Background(), "place-photos", imageName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", ErrMinioUnavailable
	}
	objInfo, err := obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, "", ErrNotFound
		}
		return nil, "", ErrMinioUnavailable
	}
	return obj, objInfo.ETag, nil
}
