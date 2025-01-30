package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/cookies"
	"forum/internal/models"
	"io/ioutil"
	"net/http"
)

const (
	githubAuthURL     = "https://github.com/login/oauth/authorize"
	githubTokenURL    = "https://github.com/login/oauth/access_token"
	githubUserInfo    = "https://api.github.com/user"
	githubRedirectURI = "https://localhost:8080/auth/github/callback"
)

type githubUser struct {
	Name   string `json:"login"`
	NodeID string `json:"node_id"`
}

// GithubLogin redirects the user to Github's OAuth authorization page
func (h *Handler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email", githubAuthURL, h.GithubConfig.ClientID, githubRedirectURI)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GithubCallback handles the callback from Github OAuth
func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	requestBodyMap := map[string]string{
		"client_id":     h.GithubConfig.ClientID,
		"client_secret": h.GithubConfig.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		h.ServerError(w, reqerr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		h.ServerError(w, resperr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	req, err := http.NewRequest(http.MethodGet, githubUserInfo, nil)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+ghresp.AccessToken)
	req.Header.Set("Accept", "application/json")

	req, reqerr = http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		h.ServerError(w, reqerr)
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", ghresp.AccessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr = http.DefaultClient.Do(req)
	if resperr != nil {
		h.ServerError(w, resperr)
	}

	respbody, _ = ioutil.ReadAll(resp.Body)

	githubUser := githubUser{}
	err = json.Unmarshal(respbody, &githubUser)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	user, err := h.service.GetUserByEmail(githubUser.NodeID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if user == nil {
		form := &models.UserSignupForm{
			Name:     githubUser.Name,
			Email:    githubUser.NodeID,
			Password: githubUser.NodeID,
		}

		err = h.service.InsertUser(form.Name, form.Email, form.Password)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	form := &models.UserLoginForm{
		Email:    githubUser.NodeID,
		Password: githubUser.NodeID,
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
