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
	GetDislikedPosts(user_id int) ([]*models.Post, error)
	AddDislikePost(post_id int, user_id int) error
	AddLikePost(post_id int, user_id int) error
	AddLikeComment(comment_id, user_id int) error
	AddDislikeComment(comment_id, user_id int) error
	DeletePost(post_id int) error
	IgnoreReport(report_id int) error
	DeleteCategory(id int) error
	CreateCategory(form models.CategoryCreateForm) error 
	DeleteComment(commentID int) error 
	GetCommentedPosts(userId int) ([]*models.CommentWithPost, error)
	UpdatePost(form models.PostCreateForm, data *models.TemplateData) error
	
}

type User interface {
	InsertUser(name, email, password string) error
	AuthenticateUser(form models.UserLoginForm, data *models.TemplateData) (*models.Session, *models.TemplateData, error)
	GetUser(r *http.Request) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	DeleteSession(token string) error
	IsSessionValid(token string) bool
	ReportPost(moderatorId int, postId int, reason string) error
	GetAllReports() ([]*models.Report, error)
	RequestModeratorRole(user_id int) error
	GetAllRequests() ([]*models.ModeratorRequest, error)
	DenyModeratorRequest(request_id int) error
	PromoteUserToModerator(request_id int) error 
	GetUserModeratorRequests(user_id int) ([]*models.ModeratorRequest, error)
	GetModeratorReports(user_id int) ([]*models.Report, error)
	DemoteModerator(userID int) error
	NotifyUser(userId int, postId int, notifType, message string) error
	GetUnreadNotifications(userId int) ([]*models.Notification, error)
	MarkNotificationAsRead(notificationId int) error
	GetNotifications() ([]*models.Notification, error)
	GetPostAuthor(post_id int) (int, error) 
}
