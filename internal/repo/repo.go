package repo

import "forum/internal/models"

type Repo interface {
	PostModel
	UserModel
	Session
	Category
	Reaction
}

type PostModel interface {
	InsertPost(user_id int, title string, content string) (int, error)
	GetPost(id int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	GetPostByCategory(id int) ([]*models.Post, error)
	GetCommentsForPost(post_id int) ([]models.Comment, error)
	GetAllComments() ([]*models.Comment, error)
	InsertComment(post_id int, user_id int, content string) error
	GetCreatedPosts(user_id int) ([]*models.Post, error)
	GetLikedPosts(user_id int) ([]*models.Post, error)
}

type UserModel interface {
	AuthenticateUser(email, password string) (int, error)
	InsertUser(name, email, password string) error
	Exists(id int) (bool, error)
	GetUserByID(id int) (*models.User, error)
}

type Session interface {
	DeleteSessionById(userId int) error
	CreateSession(session *models.Session) error
	DeleteSessionByToken(token string) error
	GetUserIDByToken(token string) (int, error)
	IsSessionValid(token string) bool
}

type Category interface {
	GetCategories() ([]*models.Category, error)
	PostCategoryPost(post_id int, category_id int) error
}

type Reaction interface {
	AddLikePost(post_id int, user_id int) error
	AddDislikePost(post_id int, user_id int) error
	CheckUserReactionsPost(post_id int, user_id int) (int, error)
	InsertUserReactionPost(post_id int, user_id int, reaction int) error
	RemoveUserReactionPost(post_id int, user_id int, reaction int) error
	AddLikeComment(comment_id int, user_id int) error
	AddDislikeComment(comment_id int, user_id int) error
	CheckUserReactionComment(comment_id int, user_id int) (int, error)
	InsertUserReactionComment(comment_id int, user_id int, reaction int) error
	RemoveUserReactionComment(comment_id int, user_id int, reaction int) error
}