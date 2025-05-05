package services

import (
	"backend/models"

	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (s *PostService) CreatePost(post *models.Post) error {
	return s.DB.Create(post).Error
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := s.DB.Preload("User").Preload("Comments").Preload("Likes").First(&post, id).Error
	return &post, err
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Preload("User").Preload("Comments").Preload("Likes").Find(&posts).Error
	return posts, err
}

func (s *PostService) FilterPostsByUser(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Where("user_id = ?", userID).Preload("User").Preload("Comments").Preload("Likes").Find(&posts).Error
	return posts, err
}

func (s *PostService) UpdatePost(post *models.Post) error {
	return s.DB.Save(post).Error
}

func (s *PostService) DeletePost(id uint) error {
	return s.DB.Delete(&models.Post{}, id).Error
}
