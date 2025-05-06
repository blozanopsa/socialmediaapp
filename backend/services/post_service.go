package services

import (
	"backend/models"
	"errors"

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
	err := s.DB.Preload("User").Preload("Likes").Preload("Comments.User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	for i := range post.Comments {
		post.Comments[i].UserName = post.Comments[i].User.Name
	}
	return &post, err
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Preload("User").Preload("Likes").Preload("Comments.User").Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	for i := range posts {
		for j := range posts[i].Comments {
			posts[i].Comments[j].UserName = posts[i].Comments[j].User.Name
		}
	}
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
	// First, delete comments associated with the post
	if err := s.DB.Where("post_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
		return err
	}
	// Then, delete the post itself
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

// GetUserByID retrieves a user by their ID.
func (s *PostService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// AddComment adds a comment to a post
func (s *PostService) AddComment(postID, userID uint, content string) (*models.Comment, error) {
	// Fetch the user
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Create the comment
	comment := models.Comment{
		PostID:   postID,
		UserID:   userID,
		User:     user,      // Keep user object associated
		UserName: user.Name, // Populate UserName
		Content:  content,
	}
	if err := s.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

// EditComment edits a comment's content
func (s *PostService) EditComment(postID, commentID uint, content string) (*models.Comment, error) {
	var comment models.Comment
	if err := s.DB.Preload("User").Where("id = ? AND post_id = ?", commentID, postID).First(&comment).Error; err != nil {
		return nil, err
	}
	comment.Content = content
	if err := s.DB.Save(&comment).Error; err != nil {
		return nil, err
	}
	// Repopulate UserName after save, in case User association was lost or not there
	if comment.User.ID != 0 { // Check if User is populated
		comment.UserName = comment.User.Name
	} else if comment.UserID != 0 { // Fallback to fetch user if User struct is empty but UserID exists
		var user models.User
		if err := s.DB.First(&user, comment.UserID).Error; err == nil {
			comment.User = user
			comment.UserName = user.Name
		}
	}
	return &comment, nil
}

// DeleteComment deletes a comment
func (s *PostService) DeleteComment(postID, commentID uint) error {
	return s.DB.Where("id = ? AND post_id = ?", commentID, postID).Delete(&models.Comment{}).Error
}
