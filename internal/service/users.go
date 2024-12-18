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
