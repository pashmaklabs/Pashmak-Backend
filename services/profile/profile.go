package services_profile

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	webp "github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models/auth"
	models "pashmak.com/pashmak/models/openai"
	models_place "pashmak.com/pashmak/models/place"
	serializers_place "pashmak.com/pashmak/serializers/place"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
	services_auth "pashmak.com/pashmak/services/auth"
)

var (
	ErrInvalidFile      = errors.New("invalid file type or size")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotFound         = errors.New("avatar not found")
	ErrMinioUnavailable = errors.New("minio unavailable")
	ErrInvalidSize      = errors.New("file too large")
)

type ProfileService struct {
	DB        *gorm.DB
	Minio     *minio.Client
	AppConfig *bootstrap.AppConfig
}

func NewProfileService(db *gorm.DB, minio *minio.Client, appConfig *bootstrap.AppConfig) *ProfileService {
	return &ProfileService{
		DB:        db,
		Minio:     minio,
		AppConfig: appConfig,
	}
}

func (ps *ProfileService) GetMyProfile(id uint) (serializers_profile.CurrentProfileResponse, error) {
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", id)
	return serializers_profile.CurrentProfileResponse{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Avatar_url: user.Avatar_url,
	}, result.Error
}

func (ps *ProfileService) GetUserByGmail(email string) (models_auth.User, error) {
	var user models_auth.User
	result := ps.DB.First(&user, "email = ?", email)
	return user, result.Error
}

func (ps *ProfileService) GetProfileByID(id uint) (serializers_profile.GetProfileByIDResponse, error) {
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", id)
	if result.Error != nil {

	}
	return serializers_profile.GetProfileByIDResponse{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Avatar_url: user.Avatar_url,
	}, result.Error
}

// CreatePlaceLabel creates a new place label for a user
func (ps *ProfileService) CreatePlaceLabel(name string, userID uint) (*models_place.PlaceLabel, error) {
	label := &models_place.PlaceLabel{
		Name:   name,
		UserID: userID,
	}

	if err := ps.DB.Create(label).Error; err != nil {
		return nil, err
	}

	return label, nil
}

// GetUserPlaceLabels returns all place labels for a specific user with the count of saved locations
func (ps *ProfileService) GetUserPlaceLabels(userID uint) ([]serializers_profile.PlaceLabelWithCountResponse, error) {
	var results []serializers_profile.PlaceLabelWithCountResponse

	if err := ps.DB.Table("place_labels").
		Select("place_labels.id, place_labels.name, COUNT(saved_locations.id) as saved_locations_count").
		Joins("LEFT JOIN saved_locations ON saved_locations.place_label_id = place_labels.id").
		Where("place_labels.user_id = ?", userID).
		Group("place_labels.id, place_labels.name").
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// DeletePlaceLabel permanently deletes a place label
func (ps *ProfileService) DeletePlaceLabel(labelID uint, userID uint) error {
	// Using Unscoped() for hard delete and checking user ownership
	result := ps.DB.Unscoped().Where("id = ? AND user_id = ?", labelID, userID).Delete(&models_place.PlaceLabel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type CreateSavedLocationParams struct {
	Latitude     float64
	Longitude    float64
	PlaceID      *uint
	PlaceLabelID uint
	UserNote     *string
	Name         *string
}

func (ps *ProfileService) CreateSavedLocation(userID uint, params CreateSavedLocationParams) (*models_place.SavedLocation, error) {
	var placeLabel models_place.PlaceLabel
	if err := ps.DB.First(&placeLabel, "id = ? AND user_id = ?", params.PlaceLabelID, userID).Error; err != nil {
		return nil, err
	}
	savedLocation := models_place.SavedLocation{
		Latitude:     params.Latitude,
		Longitude:    params.Longitude,
		PlaceLabelID: params.PlaceLabelID,
	}

	if params.PlaceID != nil {
		savedLocation.PlaceID = params.PlaceID
	}
	if params.UserNote != nil {
		savedLocation.UserNote = *params.UserNote
	}
	if params.Name != nil {
		savedLocation.Name = *params.Name
	}
	if err := ps.DB.Create(&savedLocation).Error; err != nil {
		return nil, err
	}
	return &savedLocation, nil
}

func (ps *ProfileService) GetSavedLocationsByPlaceLabel(userID uint, labelID uint) ([]models_place.SavedLocation, error) {
	var places []models_place.SavedLocation
	fmt.Println(labelID, userID)
	if err := ps.DB.
		Joins("JOIN place_labels ON place_labels.id = saved_locations.place_label_id").
		Where("saved_locations.place_label_id = ? AND place_labels.user_id = ?", labelID, userID).
		Find(&places).
		Error; err != nil {
		return nil, err
	}
	return places, nil
}

type UpdateSavedLocationServiceParams struct {
	ID           uint
	UserNote     *string
	PlaceLabelID *uint
	Name         *string
	UserID       uint
}

func (ps *ProfileService) UpdateSavedLocation(in UpdateSavedLocationServiceParams) (*models_place.SavedLocation, error) {
	var savedLocation models_place.SavedLocation
	if err := ps.DB.Joins("JOIN place_labels ON place_labels.id = saved_locations.place_label_id").Where("saved_locations.id = ? AND place_labels.user_id = ?", in.ID, in.UserID).Find(&savedLocation).Error; err != nil {
		return nil, err
	}
	if in.UserNote != nil {
		savedLocation.UserNote = *in.UserNote
	}
	if in.PlaceLabelID != nil {
		savedLocation.PlaceLabelID = *in.PlaceLabelID
	}
	if in.Name != nil {
		savedLocation.Name = *in.Name
	}
	if err := ps.DB.Save(savedLocation).Error; err != nil {
		return nil, err
	}
	return &savedLocation, nil
}

func (ps *ProfileService) HardDeleteSavedLocation(ID uint, userID uint) error {
	sub := ps.DB.
		Model(&models_place.PlaceLabel{}).
		Select("id").
		Where("user_id = ?", userID)

	result := ps.DB.
		Unscoped().
		Where("id = ? AND place_label_id IN (?)", ID, sub).
		Delete(&models_place.SavedLocation{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ps *ProfileService) GetLabelOfPlace(userID uint, placeID uint) (*serializers_place.SavedLocationResponse, error) {
	var savedLocation models_place.SavedLocation
	if err := ps.DB.
		Joins("JOIN place_labels ON place_labels.id = saved_locations.place_label_id").
		Where("place_labels.user_id = ? AND saved_locations.place_id = ?", userID, placeID).
		First(&savedLocation).
		Error; err != nil {
		return nil, err
	}
	return &serializers_place.SavedLocationResponse{
		ID:           int64(savedLocation.ID),
		PlaceLabelID: int64(savedLocation.PlaceLabelID),
	}, nil
}

func (ps *ProfileService) validateImage(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", ErrInvalidFile
	}

	if file.Size > 1<<24 {
		return "", ErrInvalidSize
	}
	return ext, nil
}

func (ps *ProfileService) uploadImage(file *multipart.FileHeader, user models_auth.User) (string, minio.UploadInfo, error) {
	var buf bytes.Buffer
	fileReader, err := file.Open()
	if err != nil {
		return "", minio.UploadInfo{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	img, _, err := image.Decode(fileReader)
	if err != nil {
		return "", minio.UploadInfo{}, err
	}

	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 30}); err != nil {
		return "", minio.UploadInfo{}, err
	}
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ".webp")
	Reader := bytes.NewReader(buf.Bytes())
	timedCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := ps.Minio.PutObject(
		timedCtx,
		"profile-photos",
		objectName,
		Reader,
		Reader.Size(),
		minio.PutObjectOptions{ContentType: "image/" + "webp"},
	)
	if err != nil {
		return "", minio.UploadInfo{}, err
	}

	user.Avatar_url = fmt.Sprintf("%s/profiles/avatar/%s", ps.AppConfig.ServerHost, objectName)
	saveres := ps.DB.Save(user)
	if saveres.Error != nil {
		return "", minio.UploadInfo{}, saveres.Error
	}

	return objectName, info, nil
}

func (ps *ProfileService) GetAvatar(ctx context.Context, fileName string, height int) (io.ReadCloser, string, error) {
	if fileName == "" {
		return nil, "", ErrInvalidFile
	}

	obj, err := ps.Minio.GetObject(ctx, "profile-photos", fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", ErrMinioUnavailable
	}

	// Check if object exists by trying to get stats
	objInfo, err := obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, "", ErrNotFound
		}
		return nil, "", ErrMinioUnavailable
	}

	// If no resizing requested, return the response body directly
	if height == 0 {
		return obj, objInfo.ETag, nil
	}

	// Read and resize the image
	data, err := io.ReadAll(obj)
	obj.Close()
	if err != nil {
		return nil, "", err
	}

	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	resized := imaging.Resize(img, height, 0, imaging.Lanczos)
	buf := new(bytes.Buffer)
	if err = webp.Encode(buf, resized, &webp.Options{Lossless: false, Quality: 30}); err != nil {
		return nil, "", err
	}

	return io.NopCloser(buf), objInfo.ETag, nil
}

func (ps *ProfileService) UploadUserAvatar(ctx *gin.Context, userID uint) (resp gin.H, err error) {
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return gin.H{}, ErrNotFound
		}
		return gin.H{}, result.Error
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		return nil, fmt.Errorf("failed to get file from form: %w", err)
	}

	_, err = ps.validateImage(file)
	if err != nil {
		return nil, fmt.Errorf("failed to validate image: %w", err)
	}

	objectName, info, err := ps.uploadImage(file, user)
	if err != nil {
		return nil, fmt.Errorf("failed to put object to minio: %w", err)
	}

	return gin.H{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"objectName": objectName,
			"info":       info,
		},
	}, nil

	// TODO: Remove old avatar if exists
	// TODO: Validate MIME Type
	// TODO: Sanitize File Content: Use an image processing library (e.g., imaging or bimg) to validate and sanitize the image, ensuring it's a valid image and not malicious.
	// TODO: Resize the image to a standard size (e.g., 256x256 pixels) using an image processing library.
	// TODO: Asynchronous Compression
	// TODO: Some suggestions: https://x.com/i/grok/share/LOi6Xexr8xBCaX49J7t0LrUgN
	// TODO: Rate Limiting
	// TODO: Store Thumbnails, Medium, Full-size
}

func (ps *ProfileService) UpdateUserProfile(userInfo services_auth.UserInfo, payload serializers_profile.UpdateUserProfileRequest) error {
	user, err := ps.GetUserByGmail(userInfo.Email)
	if err != nil {
		return err
	}
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName

	if err := ps.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ps *ProfileService) FetchSearchHistory(userInfo services_auth.UserInfo) ([]models.SearchHistory, error) {
	var history []models.SearchHistory
	historyQuery := ps.DB.
		Where("user_id = ?", userInfo.ID).
		Find(&history)

	if historyQuery.Error != nil {
		return nil, historyQuery.Error
	}

	if len(history) == 0 {
		return nil, errors.New("no history found")
	}
	return history, nil
}

func (ps *ProfileService) DeleteSearchHistory(userInfo services_auth.UserInfo, searchId string) error {
	var history models.SearchHistory
	if err := ps.DB.First(&history, searchId).Error; err != nil {
		return errors.New("history not found")
	}

	if err := ps.DB.Where("id = ? AND user_id = ?", searchId, userInfo.ID).Delete(&history).Error; err != nil {
		return err
	}
	return nil
}

func (ps *ProfileService) ClearSearchHistory(userInfo services_auth.UserInfo) error {
	var history models.SearchHistory
	if err := ps.DB.Where("user_id = ?", userInfo.ID).Delete(&history).Error; err != nil {
		return err
	}
	return nil
}
