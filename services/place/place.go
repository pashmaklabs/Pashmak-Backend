package services_place

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	webp "github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	oa "pashmak.com/pashmak/models/openai"
	models_place "pashmak.com/pashmak/models/place"
	serializers_place "pashmak.com/pashmak/serializers/place"
	sp "pashmak.com/pashmak/serializers/place"
	services_auth "pashmak.com/pashmak/services/auth"
	services_openai "pashmak.com/pashmak/services/openai"
	"pashmak.com/pashmak/services/placeOsmUtils"
)

var (
	ErrInvalidFile      = errors.New("invalid file type or size")
	ErrNotFound         = errors.New("image not found")
	ErrMinioUnavailable = errors.New("minio unavailable")
	ErrInvalidSize      = errors.New("file too large")
)

type placeSearchResult struct {
	OsmID     *uint
	Name      string
	Amenity   *string
	Latitude  *float64
	Longitude *float64
	ID        *int64
}

type PlaceService struct {
	DB            *gorm.DB
	AppConfig     *bootstrap.AppConfig
	OpenAIService *services_openai.OpenAIService
	Minio         *minio.Client
	PGVectorDB    *gorm.DB
}

func NewPlaceService(db *gorm.DB, appConfig *bootstrap.AppConfig, openaiService *services_openai.OpenAIService, minioClient *minio.Client, pgvectorDB *gorm.DB) *PlaceService {
	return &PlaceService{
		DB:            db,
		AppConfig:     appConfig,
		OpenAIService: openaiService,
		Minio:         minioClient,
		PGVectorDB:    pgvectorDB,
	}
}

func (ps *PlaceService) GetPlaceByID(id string) (*sp.GetPlaceByIDResponse, error) {
	// Try to parse the string ID to uint first
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		// If parsing fails, search in gplace table
		var result struct {
			ID        uint    `db:"id"`
			Name      string  `db:"name"`
			GmapID    string  `db:"gmap_id"`
			Category  string  `db:"category"`
			Latitude  float64 `db:"latitude"`
			Longitude float64 `db:"longitude"`
		}

		if err := ps.PGVectorDB.Raw("SELECT id, name, gmap_id, category, latitude, longitude FROM gplaces WHERE gmap_id = ? AND deleted_at IS NULL", id).First(&result).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("no place found with ID %s", id)
			}
			return nil, fmt.Errorf("error searching gplace table: %w", err)
		}

		// Debug: print the category value
		fmt.Printf("DEBUG: Category value: %+v, Type: %T\n", result.Category, result.Category)

		// Parse PostgreSQL array string format (e.g., "{Restaurant}" -> ["Restaurant"])
		amenity := ""
		if result.Category != "" {
			// Remove the curly braces and split by comma
			categoryStr := strings.Trim(result.Category, "{}")
			if categoryStr != "" {
				categories := strings.Split(categoryStr, ",")
				if len(categories) > 0 {
					amenity = strings.TrimSpace(categories[0])
				}
			}
		}

		// Create response from gplace data
		response := sp.GetPlaceByIDResponse{
			ID:        result.GmapID, // Use GmapID as the string ID
			Name:      result.Name,
			Amenity:   &amenity, // Use first category as amenity
			Latitude:  &result.Latitude,
			Longitude: &result.Longitude,
			ImageURLs: []string{}, // gplace doesn't have images in our current structure
		}

		return &response, nil
	}

	placeID := uint(idUint)

	// First, try to import/get the place from OSM if it exists
	place, err := placeOsmUtils.ImportFromOSM(placeID, ps.DB)
	if err != nil {
		return nil, err
	}

	// Load the place with its images
	res := ps.DB.Preload("Images").First(&place, placeID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no place found with ID %s", id)
		}
		return nil, res.Error
	}

	var response sp.GetPlaceByIDResponse

	// If this is an OSM place, get additional data from planet_osm_point
	if place.IsOSM {
		var osmResult struct {
			OsmID     uint    `db:"osm_id"`
			Name      string  `db:"name"`
			Amenity   string  `db:"amenity"`
			Latitude  float64 `db:"latitude"`
			Longitude float64 `db:"longitude"`
		}

		query := `
			SELECT
				osm_id,
				name,
				amenity,
				ST_Y(ST_Transform(way, 4326)) as latitude,
				ST_X(ST_Transform(way, 4326)) as longitude
			FROM planet_osm_point
			WHERE osm_id = ?`

		err = ps.DB.Raw(query, placeID).Scan(&osmResult).Error
		if err != nil {
			return nil, err
		}

		// Populate response with OSM data
		response.OsmID = &osmResult.OsmID
		response.Name = osmResult.Name
		response.Amenity = &osmResult.Amenity
		response.Latitude = &osmResult.Latitude
		response.Longitude = &osmResult.Longitude
	} else {
		// For non-OSM places, use data from the places table
		response.OsmID = nil
		response.Name = place.Name
		response.Amenity = &place.Amenity
		response.Latitude = &place.Latitude
		response.Longitude = &place.Longitude
	}

	// Common data for both OSM and non-OSM places
	response.ID = strconv.FormatUint(uint64(place.ID), 10) // Convert uint to string

	// Add image URLs
	for _, image := range place.Images {
		response.ImageURLs = append(response.ImageURLs, image.URL)
	}

	return &response, nil
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

	if err := ps.DB.Create(&history); err != nil {
		return err.Error
	}
	return nil
}

func (ps *PlaceService) SearchPlace(q string, lat string, long string, agentic bool) ([]sp.GetPlaceByIDResponse, error) {
	query := fmt.Sprintf("Query: %s\nLatitude: %s\nLongitude: %s", q, lat, long)
	var rawResults []placeSearchResult

	if agentic {
		categories, categoryErr := ps.GetQueryCategories(q)
		if categoryErr != nil {
			return nil, categoryErr
		}
		var results []sp.GetPlaceByIDResponse
		vectorSearchTool := NewVectorSearchTool(ps.PGVectorDB, "greviews", ps.AppConfig)

		// Create proper JSON for categories
		categoriesJSON, _ := json.Marshal(categories)
		rawString, _ := vectorSearchTool.Call(context.Background(), fmt.Sprintf(`{"query": "%s", "limit": 20, "categories": %s}`, q, string(categoriesJSON)))
		rawResults := []placeRow{}
		fmt.Println("rawString", rawString)
		err := json.Unmarshal([]byte(rawString), &rawResults)
		if err != nil {
			return nil, err
		}
		fmt.Println("rawResults", rawResults)
		for _, result := range rawResults {
			fmt.Println("result", result.GMapID)
			place, err := ps.GetPlaceByID(result.GMapID)
			if err != nil {
				return nil, err
			}
			results = append(results, *place)
		}
		return results, nil
	}

	// Create a new SQL chat agent with our specialized system prompt
	sqlAgent := oa.NewSQLChatAgent("", "gpt-4.1-mini")
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
			SELECT NULL AS osm_id, name, NULL AS latitude, NULL AS longitude, NULL as amenity, id
			FROM places
			WHERE name ILIKE ?
			UNION
			SELECT osm_id, name, ST_Y(ST_Transform(way, 4326)) AS latitude, ST_X(ST_Transform(way, 4326)) AS longitude, amenity, NULL AS id
			FROM planet_osm_point
			WHERE name ILIKE ?
			LIMIT 50`

		err = ps.DB.Raw(fallbackQuery, "%"+q+"%", "%"+q+"%").Scan(&rawResults).Error
		if err != nil {
			log.Printf("Fallback search failed: %v", err)
			return nil, fmt.Errorf("fallback search failed: %w", err)
		}
		return ps.mapSearchResultsToResponse(rawResults), nil
	}

	if len(resp.Choices) == 0 {
		log.Printf("No SQL query generated from OpenAI response: %+v", resp)
		return nil, fmt.Errorf("no SQL query generated")
	}

	// Execute the generated SQL query
	generatedSQL := resp.Choices[0].Message.Content
	log.Printf("Generated SQL: %s", generatedSQL)

	err = ps.DB.Raw(generatedSQL).Scan(&rawResults).Error
	if err != nil {
		log.Printf("Failed to execute generated SQL: %v\nSQL: %s", err, generatedSQL)
		return nil, fmt.Errorf("failed to execute generated SQL: %w", err)
	}

	return ps.mapSearchResultsToResponse(rawResults), nil
}

// Helper function to map search results to response
func (ps *PlaceService) mapSearchResultsToResponse(rawResults []placeSearchResult) []sp.GetPlaceByIDResponse {
	var results []sp.GetPlaceByIDResponse
	for _, r := range rawResults {
		resp := sp.GetPlaceByIDResponse{
			OsmID:     r.OsmID,
			Name:      r.Name,
			Amenity:   r.Amenity,
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
			ID:        "",
			ImageURLs: []string{},
		}
		if r.ID != nil {
			resp.ID = strconv.FormatInt(*r.ID, 10)
			// Optionally, fetch images for this place ID
			var images []models_place.Image
			if err := ps.DB.Where("place_id = ?", *r.ID).Find(&images).Error; err == nil {
				for _, img := range images {
					resp.ImageURLs = append(resp.ImageURLs, img.URL)
				}
			}
		}
		results = append(results, resp)
	}
	return results
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
func (ps *PlaceService) UploadPlaceImage(place *models_place.Place, file *multipart.FileHeader) (string, error) {
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
	imgURL := fmt.Sprintf("%s/places/%d/images/%s", ps.AppConfig.ServerHost, place.ID, objectName)
	res := ps.DB.Create(&models_place.Image{
		PlaceID: place.ID,
		URL:     imgURL,
		AltText: "alt",
		Caption: "caption",
	})

	if res.Error != nil {
		return "", res.Error
	}

	return imgURL, nil
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

func (ps *PlaceService) AddNewPlace(userinfo services_auth.UserInfo, payload serializers_place.AddPlaceRequest) error {
	newPlace := models_place.Place{
		Name:      payload.Name,
		Amenity:   payload.Amenity,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		IsOSM:     false,
	}
	if err := ps.DB.Create(&newPlace).Error; err != nil {
		return err
	}
	return nil
}

type VectorSearchTool struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
	TableName string
}

func NewVectorSearchTool(db *gorm.DB, tableName string, appConfig *bootstrap.AppConfig) *VectorSearchTool {
	return &VectorSearchTool{DB: db, TableName: tableName, AppConfig: appConfig}
}

func (t *VectorSearchTool) Name() string {
	return "vector_search"
}

func (t *VectorSearchTool) Description() string {
	return "Return up to `limit` places most relevant to `query`, as JSON list of {gmap_id, name, review, rating}"
}

type placeRow struct {
	GMapID     string  `gorm:"column:gmap_id" json:"gmap_id"`
	Name       string  `gorm:"column:name"    json:"name"`
	Review     string  `gorm:"column:review"  json:"review"`
	Rating     float64 `gorm:"column:rating"  json:"rating"`
	Similarity float64 `gorm:"column:similarity"  json:"similarity"`
}

// Call is invoked by the LLM (via function‑calling). argsJSON looks like: {"query":"…","limit":5}
func (t *VectorSearchTool) Call(ctx context.Context, argsJSON string) (string, error) {
	fmt.Println("argsJSON", argsJSON)
	// 1️⃣ parse the arguments
	var args struct {
		Query      string   `json:"query"`
		Limit      int      `json:"limit"`
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("vector_search: bad args: %w", err)
	}

	category_filter := ""
	if len(args.Categories) > 0 {
		category_filter = `
			AND p.category && $4
		`
	}
	var rows []placeRow
	sql := fmt.Sprintf(`
        		SELECT 
                    r.text as review,
                    r.rating as rating,
                    1 - (r.embedding <=> $1::vector) as similarity,
                    p.gmap_id as gmap_id,
					p.name as name
                FROM %s r
                LEFT JOIN gplaces p ON r.gmap_id = p.gmap_id
                WHERE r.embedding IS NOT NULL
                %s
                ORDER BY r.embedding <=> $2::vector
                LIMIT $3`, t.TableName, category_filter)

	// for simplicity we assume you've already got queryEmbedding as JSON array text
	// you could factor out your getEmbedding() if you like
	// here we inline it:
	emb, err := t.getEmbedding(ctx, args.Query)
	if err != nil {
		return "", err
	}
	embArg, _ := json.Marshal(emb) // "[0.12,0.53,…]"

	var err2 error
	if len(args.Categories) > 0 {
		// Execute with category filter
		err2 = t.DB.Raw(sql, string(embArg), string(embArg), args.Limit, pq.Array(args.Categories)).Scan(&rows).Error
	} else {
		// Execute without category filter
		err2 = t.DB.Raw(sql, string(embArg), string(embArg), args.Limit).Scan(&rows).Error
	}

	if err2 != nil {
		return "", fmt.Errorf("vector_search query: %w", err2)
	}

	// Make results unique by gmap_id
	uniqueRows := make([]placeRow, 0)
	seenGMapIDs := make(map[string]bool)

	for _, row := range rows {
		if !seenGMapIDs[row.GMapID] {
			seenGMapIDs[row.GMapID] = true
			uniqueRows = append(uniqueRows, row)
		}
	}

	// 3️⃣ marshal the result list to JSON
	out, err := json.Marshal(uniqueRows)
	if err != nil {
		return "", fmt.Errorf("vector_search marshal: %w", err)
	}

	return string(out), nil
}

// getEmbedding can be unexported inside the tool
func (t *VectorSearchTool) getEmbedding(ctx context.Context, text string) ([]float32, error) {
	client := openai.NewClient(option.WithAPIKey(t.AppConfig.OpenaiApiKey))
	resp, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Small,
		Input: openai.EmbeddingNewParamsInputUnion{OfString: openai.String(text)},
	})
	if err != nil {
		return nil, fmt.Errorf("getEmbedding: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("getEmbedding: no data")
	}
	emb := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		emb[i] = float32(v)
	}
	return emb, nil
}

func (ps *PlaceService) GetPlaceRecommendations(query string) ([]string, error) {
	fmt.Println(query)
	ctx := context.Background()
	// prepare the tool in the SDK format
	vectorToolParam := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "vector_search",
			Description: openai.String("Return up to `limit` places most relevant to `query`, as JSON list of {gmap_id, name, review, rating}"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{"type": "string", "description": "Descriptive place query"},
					"limit": map[string]interface{}{"type": "integer", "description": "How many results to return"},
				},
				"required": []string{"query"},
			},
		},
	}

	client := openai.NewClient(
		option.WithAPIKey(ps.AppConfig.OpenaiApiKey),
	)

	// seed the chat history
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
		You are a place‑recommendation assistant. Use the vector_search tool if you need to look up places by semantic similarity.
		You should return a list of places that are most relevant to the query.
		Don't ask questions, just return the list of places.
		All I want you is a list of gmap_ids.
		`),
		openai.UserMessage(query),
	}

	const maxRounds = 3
	for round := 0; round < maxRounds; round++ {
		// 1) ask GPT (it may produce a tool call)
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModelGPT4o,
			Messages: messages,
			Tools:    []openai.ChatCompletionToolParam{vectorToolParam},
		})
		if err != nil {
			return nil, err
		}

		msg := resp.Choices[0].Message
		// fmt.Println("msg", msg)
		// 2) did it call our tool?
		if len(msg.ToolCalls) > 0 {
			call := msg.ToolCalls[0]
			// fmt.Println("call.Function.Name", call.Function.Name)
			if call.Function.Name == vectorToolParam.Function.Name {
				// execute the tool with the JSON args
				// fmt.Println("call.Function.Arguments", call.Function.Arguments)
				toolResult, err := NewVectorSearchTool(ps.PGVectorDB, "greviews", ps.AppConfig).Call(ctx, call.Function.Arguments)
				if err != nil {
					return nil, err
				}
				fmt.Println("toolResult", toolResult)
				// append the assistant’s previous message *and* the tool response
				messages = append(messages,
					msg.ToParam(),                           // the function_call message
					openai.ToolMessage(toolResult, call.ID), // the tool’s JSON reply
				)
				// loop again so GPT can pick the best gmap_id
				continue
			}
		}

		// 3) no tool call → this is the final answer
		return []string{msg.Content}, nil
	}

	return []string{}, fmt.Errorf("GPT never returned a final answer after %d rounds", maxRounds)
}

func (ps *PlaceService) GetQueryCategories(query string) ([]string, error) {
	categories := []string{
		"Restaurant",
		"Coffee shop",
		"Bar",
		"Park",
		"Tourist attraction",
		"Clothing store",
		"Gift shop",
		"Auto repair shop",
		"Hotel",
		"Condominium complex",
		"Medical clinic",
		"Grocery store",
		"Electronics store",
		"Home goods store",
		"Beauty salon",
		"Bank",
		"Gas station",
		"Pharmacy",
		"Car rental agency",
		"Movie theater",
	}

	// Create OpenAI client
	client := openai.NewClient(option.WithAPIKey(ps.AppConfig.OpenaiApiKey))

	// Create system prompt for category detection
	systemPrompt := fmt.Sprintf(`You are a place category classifier. Given a user query, return one or more relevant categories from this exact list:

%s

Rules:
1. Only return categories that are EXACTLY in the list above
2. Return multiple categories if the query could match multiple types of places
3. If no category matches, return an empty list
4. Return only the category names, separated by commas
5. Do not add any explanations or additional text

Example queries and responses:
- "I want to eat pizza" → "Restaurant"
- "Looking for coffee and wifi" → "Coffee shop"
- "Need to buy clothes and get a haircut" → "Clothing store,Beauty salon"
- "Where can I park my car?" → "Park"
- "Random unrelated query" → ""`, strings.Join(categories, "\n"))

	// Create messages for the chat completion
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(query),
	}

	// Call OpenAI API
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4oMini,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return []string{}, nil
	}

	// Extract the response content
	response := strings.TrimSpace(resp.Choices[0].Message.Content)

	// If response is empty, return empty slice
	if response == "" {
		return []string{}, nil
	}

	// Split by comma and clean up each category
	categoryList := strings.Split(response, ",")
	var result []string

	for _, cat := range categoryList {
		cat = strings.TrimSpace(cat)
		if cat != "" {
			// Verify the category exists in our predefined list
			for _, validCat := range categories {
				if cat == validCat {
					result = append(result, cat)
					break
				}
			}
		}
	}

	return result, nil
}
