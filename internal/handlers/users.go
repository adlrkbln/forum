package handlers

import (
	"errors"
	"forum/internal/cookies"
	"forum/internal/models"
	"forum/internal/validate"
	"net/http"
	"strconv"
)

func (h *Handler) userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = models.UserSignupForm{}
		h.Render(w, http.StatusOK, "signup.tmpl", data)
	case http.MethodPost:
		h.userSignupPost(w, r)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserSignupForm

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form.Name = r.PostForm.Get("name")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	form.CheckField(validate.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validate.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validate.Matches(form.Email, validate.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validate.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validate.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = h.service.InsertUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data, err := h.NewTemplateData(w, r)
			if err != nil {
				h.ServerError(w, err)
				return
			}
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			h.ServerError(w, err)
		}

		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handler) userLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = models.UserLoginForm{}
		h.Render(w, http.StatusOK, "login.tmpl", data)
	case http.MethodPost:
		h.userLoginPost(w, r)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form models.UserLoginForm

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	form.CheckField(validate.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validate.Matches(form.Email, validate.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validate.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	session, data, err := h.service.AuthenticateUser(form, data)
	if err != nil {
		if err == models.ErrNotValidPostForm {
			h.Render(w, http.StatusBadRequest, "login.tmpl", data)
			return
		}
		h.ServerError(w, err)
		return
	}
	cookies.SetSessionCookie("session_id", w, session.Token, session.ExpTime)
	http.Redirect(w, r, "/post/create", http.StatusSeeOther)
}

func (h *Handler) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	c := cookies.GetSessionCookie("session_id", r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		cookies.ExpireSessionCookie("session_id", w)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	comments, err := h.service.GetAllComments()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("CommentId"))
	if err != nil || commentID <= 0 || !CommentExists(commentID, comments) {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	if user.Role != "Admin" {
		http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteComment(commentID); err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}