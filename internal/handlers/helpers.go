package handlers

import (
	"bytes"
	"fmt"
	"forum/internal/app"
	"forum/internal/cookies"
	"forum/internal/models"
	"forum/internal/service"
	"net/http"
	"time"
)

type Handler struct {
	*app.Application
	service service.Service
}

func New(a *app.Application, r service.Service) *Handler {
	return &Handler{
		a,
		r,
	}
}

func (h *Handler) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, check := h.TemplateCache[page]
	if !check {
		err := fmt.Errorf("the template %s does not exist", page)
		h.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (h *Handler) NewTemplateData(w http.ResponseWriter, r *http.Request) (*models.TemplateData, error) {
	var TemplateData models.TemplateData

	TemplateData.IsAuthenticated = h.IsAuthenticated(r)

	if TemplateData.IsAuthenticated {
		user, err := h.service.GetUser(r)
		if err != nil {
			// If the session is invalid (e.g., no rows in the database), expire the cookie
			if err == models.ErrInvalidSession {
				cookies.ExpireSessionCookie("session_id", w)
				TemplateData.IsAuthenticated = false
			} else {
				return nil, err
			}
		} else {
			TemplateData.User = user
		}
	}

	TemplateData.CurrentYear = time.Now().Year()
	return &TemplateData, nil
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	cookie := cookies.GetSessionCookie("session_id", r)
	return cookie != nil && cookie.Value != ""
}
