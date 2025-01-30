package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/cookies"
	"forum/internal/models"
	"net/http"
	"net/url"
	"strings"
)

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	googleUserInfo = "https://www.googleapis.com/oauth2/v2/userinfo"
	redirectURI    = "https://localhost:8080/auth/google/callback"
	scope          = "email profile"
)

type tokenResp struct {
	AccessToken string `json:"access_token"`
}

type googleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}

// GoogleLogin redirects the user to Google's OAuth 2.0 authorization page
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&access_type=offline", googleAuthURL, h.GoogleConfig.ClientID, url.QueryEscape(redirectURI), url.QueryEscape(scope))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GoogleCallback handles the callback from Google OAuth
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	data_url := url.Values{}
	data_url.Set("code", code)
	data_url.Set("client_id", h.GoogleConfig.ClientID)
	data_url.Set("client_secret", h.GoogleConfig.ClientSecret)
	data_url.Set("redirect_uri", redirectURI)
	data_url.Set("grant_type", "authorization_code")

	resp, err := http.Post(googleTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data_url.Encode()))
	if err != nil {
		h.ServerError(w, err)
		return
	}
	defer resp.Body.Close()

	tokenResp := &tokenResp{}

	err = json.NewDecoder(resp.Body).Decode(tokenResp)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, googleUserInfo, nil)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	defer resp.Body.Close()

	googleUser := googleUser{}
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	user, err := h.service.GetUserByEmail(googleUser.Email)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if user == nil {
		form := &models.UserSignupForm{
			Name:     googleUser.Name,
			Email:    googleUser.Email,
			Password: googleUser.Sub,
		}

		err = h.service.InsertUser(form.Name, form.Email, form.Password)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	form := &models.UserLoginForm{
		Email:    googleUser.Email,
		Password: googleUser.Sub,
	}
	data, err := h.NewTemplateData(w, r)
	session, data, err := h.service.AuthenticateUser(*form, data)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	cookies.SetSessionCookie("session_id", w, session.Token, session.ExpTime)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
