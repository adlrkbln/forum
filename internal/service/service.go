package service

import (
	"forum/internal/models"
	"forum/internal/repo"
	"net/http"
)

type service struct {
	repo repo.Repo
}

type Service interface {
	Post
	User
}

func NewService(repo repo.Repo) Service {
	return &service{repo}
}

type Post interface {
	InsertPost(form models.PostCreateForm, data *models.TemplateData) (int, error)
	GetPost(id int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	GetCategories() ([]*models.Category, error)
	GetPostByCategory(id int) ([]*models.Post, error)
	PostCategoryPost(post_id int, categoryIds []int) error
	GetAllComments() ([]*models.Comment, error)
	InsertComment(post_id int, user_id int, content string) error
	GetCreatedPosts(user_id int) ([]*models.Post, error)
	GetLikedPosts(user_id int) ([]*models.Post, error)
	AddDislikePost(post_id int, user_id int) error
	AddLikePost(post_id int, user_id int) error
	AddLikeComment(comment_id, user_id int) error
	AddDislikeComment(comment_id, user_id int) error
}

type User interface {
	InsertUser(name, email, password string) error
	AuthenticateUser(form models.UserLoginForm, data *models.TemplateData) (*models.Session, *models.TemplateData, error)
	GetUser(r *http.Request) (*models.User, error)
	DeleteSession(token string) error
	IsSessionValid(token string) bool
}
