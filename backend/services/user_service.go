package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// FindOrCreateByEmail finds a user by email or creates a new one
func (s *UserService) FindOrCreateByEmail(name, email, provider string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = models.User{Name: name, Email: email, Provider: provider}
		// Do NOT set SessionID here!
		if err := s.DB.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	} else if err != nil {
		return nil, err
	}
	// Optionally update name if changed
	if user.Name != name {
		user.Name = name
		s.DB.Save(&user)
	}
	return &user, nil
}

// GetUserBySessionID fetches a user by session ID
func (s *UserService) GetUserBySessionID(sessionID string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("session_id = ?", sessionID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveSessionID saves the session ID for a user
func (s *UserService) SaveSessionID(userID uint, sessionID string) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Update("session_id", sessionID).Error
}

// DeleteSessionID removes the session ID from the user with the given sessionID
func (s *UserService) DeleteSessionID(sessionID string) error {
	return s.DB.Model(&models.User{}).Where("session_id = ?", sessionID).Update("session_id", "").Error
}

// GetAllUsers returns all users from the database
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := s.DB.Find(&users).Error
	return users, err
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.DB.First(&user, id).Error
	return &user, err
}
