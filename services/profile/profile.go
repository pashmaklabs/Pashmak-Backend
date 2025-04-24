package services_profile

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	// "time"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models/auth"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
)

var (
    ErrInvalidFile       = errors.New("invalid file type or size")
    ErrPermissionDenied  = errors.New("permission denied")
    ErrNotFound          = errors.New("avatar not found")
    ErrMinioUnavailable  = errors.New("minio unavailable")
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

func (ps *ProfileService) GetAvatarViaPresignedURL(ctx context.Context, userID string, height int)(io.ReadCloser, string, error){
	if userID == "" {
        return nil, "", ErrInvalidFile
    }
	objectPath := userID + ".png"
	obj, err := ps.Minio.GetObject(ctx, "profile-photos", objectPath, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return nil, "", ErrMinioUnavailable
	}

	// Create a new http.Response to maintain compatibility with the rest of the function
	resp := &http.Response{
		Body:       obj,
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
	}

	// Check if object exists by trying to get stats
	_, err = obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, "", ErrNotFound
		}
		return nil, "", ErrMinioUnavailable
	}

	

	if resp.StatusCode == http.StatusNotFound {
        resp.Body.Close()
        return nil, "", ErrNotFound
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, "", ErrMinioUnavailable
    }

    eTag := resp.Header.Get("ETag")

    // If no resizing requested, return the response body directly
    if height == 0 {
        return resp.Body, eTag, nil
    }

    // Read and resize the image
    data, err := io.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return nil, "", err
    }

    img, err := imaging.Decode(bytes.NewReader(data))
    if err != nil {
        return nil, "", err
    }

    resized := imaging.Resize(img, height, 0, imaging.Lanczos)
    buf := new(bytes.Buffer)
    err = imaging.Encode(buf, resized, imaging.PNG)
    if err != nil {
        return nil, "", err
    }

    return io.NopCloser(buf), eTag, nil
}

func (ps *ProfileService) UploadUserAvatar() {

}	
