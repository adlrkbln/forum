package service

import (
	"database/sql"
	"errors"
	"forum/internal/cookies"
	"forum/internal/models"
	"net/http"
)

func (s *service) InsertUser(name, email, password string) error {
	err := s.repo.InsertUser(name, email, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AuthenticateUser(form models.UserLoginForm, data *models.TemplateData) (*models.Session, *models.TemplateData, error) {
	userId, err := s.repo.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("email", "Email or password is incorrect")
			data.Form = form
			return nil, data, models.ErrNotValidPostForm
		} else {
			return nil, nil, err
		}
	}
	session := models.NewSession(userId)
	if err = s.repo.DeleteSessionById(userId); err != nil {
		return nil, nil, err
	}
	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, nil, err
	}
	data.Form = form
	return session, data, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (s *service) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *service) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}

func (s *service) IsSessionValid(token string) bool {
	return s.repo.IsSessionValid(token)
}

func (s *service) GetUser(r *http.Request) (*models.User, error) {
	token := cookies.GetSessionCookie("session_id", r)
	if token == nil || token.Value == "" {
		return nil, models.ErrInvalidSession
	}

	userID, err := s.repo.GetUserIDByToken(token.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidSession
		}
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

func (s *service) ReportPost(moderator_id int, post_id int, reason string) error {
	return s.repo.InsertReportPost(moderator_id, post_id, reason)
}

func (s *service) GetAllReports() ([]*models.Report, error) {
	reports, err := s.repo.GetAllReports()
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *service) RequestModeratorRole(user_id int) error {
	return s.repo.RequestModeratorRole(user_id)
}

func (s *service) GetAllRequests() ([]*models.ModeratorRequest, error) {
	requests, err := s.repo.GetAllRequests()
	if err != nil {
		return nil, err
	}
	return requests, err
}

func (s *service) PromoteUserToModerator(request_id int) error {
	return s.repo.PromoteUserToModerator(request_id)
}

func (s *service) DenyModeratorRequest(request_id int) error {
	return s.repo.DenyModeratorRequest(request_id)
}

func (s *service) GetUserModeratorRequests(user_id int) ([]*models.ModeratorRequest, error) {
	requests, err := s.repo.GetUserModeratorRequests(user_id)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *service) GetModeratorReports(user_id int) ([]*models.Report, error) {
	reports, err := s.repo.GetModeratorReports(user_id)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *service) DemoteModerator(userID int) error {
	return s.repo.DemoteModerator(userID)
}

func (s *service) GetPostAuthor(post_id int) (int, error) {
	author, err := s.repo.GetPostAuthor(post_id)
	if err != nil {
		return 0, err
	}
	return author.Id, nil
}

func (s *service) GetCommentAuthor(comment_id int) (int, error) {
	author, err := s.repo.GetCommentAuthor(comment_id)
	if err != nil {
		return 0, err
	}
	return author.Id, nil
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}