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

func (ps *ProfileService) GetUserByGmail(email string) (models_auth.User, error){
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

func (ps *ProfileService) validateImage(file *multipart.FileHeader) (string, error){
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", ErrInvalidFile
	}

	if file.Size > 1<<24 {
		return "", ErrInvalidSize
	}
	return ext, nil
}

func (ps *ProfileService) uploadImage(file *multipart.FileHeader, user models_auth.User)(string, minio.UploadInfo, error){
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
	if err != nil{
		return "",  minio.UploadInfo{}, err
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

func (ps *ProfileService) UploadUserAvatar(ctx *gin.Context, userID string) (resp gin.H, err error) {
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

func (ps *ProfileService) UpdateUserProfile(userInfo services_auth.UserInfo, payload serializers_profile.UpdateUserProfileRequest) error{
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