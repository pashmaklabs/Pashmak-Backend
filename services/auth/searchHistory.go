package services_auth

import (
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
	models "pashmak.com/pashmak/models/openai"
)

func (as *AuthService) MergeSearchHistory(sessionID string, userInfo models_auth.User) error{
	err := as.DB.Transaction(func(tx *gorm.DB) error {
		// Update anonymous history to user_id
		if err := tx.Model(&models.SearchHistory{}).
			Where("session_id = ?", sessionID).
			Update("user_id", userInfo.ID).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// Clear session_id
		if err := tx.Model(&models.SearchHistory{}).
			Where("session_id = ? AND user_id = ?", sessionID, userInfo.ID).
			Update("session_id", "").Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		return nil
	})
	return err
}