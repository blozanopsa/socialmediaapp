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

// LikePost adds a like for a post by a user
func (s *PostService) LikePost(postID, userID uint) error {
	like := models.Like{PostID: postID, UserID: userID}
	return s.DB.FirstOrCreate(&like, models.Like{PostID: postID, UserID: userID}).Error
}

// UnlikePost removes a like for a post by a user
func (s *PostService) UnlikePost(postID, userID uint) error {
	return s.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.Like{}).Error
}

// GetPostsLikedByUser returns all posts liked by a specific user
func (s *PostService) GetPostsLikedByUser(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Joins("JOIN likes ON likes.post_id = posts.id").
		Where("likes.user_id = ?", userID).
		Preload("User").Preload("Comments").Preload("Likes").
		Find(&posts).Error
	return posts, err
}

// AddComment adds a comment to a post
func (s *PostService) AddComment(postID, userID uint, content string) (*models.Comment, error) {
	comment := models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	err := s.DB.Create(&comment).Error
	return &comment, err
}

// EditComment edits a comment's content
func (s *PostService) EditComment(postID, commentID uint, content string) (*models.Comment, error) {
	var comment models.Comment
	if err := s.DB.Where("id = ? AND post_id = ?", commentID, postID).First(&comment).Error; err != nil {
		return nil, err
	}
	comment.Content = content
	if err := s.DB.Save(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment deletes a comment
func (s *PostService) DeleteComment(postID, commentID uint) error {
	return s.DB.Where("id = ? AND post_id = ?", commentID, postID).Delete(&models.Comment{}).Error
}
