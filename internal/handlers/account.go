package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) accountPageGet(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data.User = user

	if user.Role == "User" {
		requests, err := h.service.GetUserModeratorRequests(user.Id)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.ModeratorRequests = requests
	}

	if user.Role == "Moderator" {
		reports, err := h.service.GetModeratorReports(user.Id)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Reports = reports
	}

	if user.Role == "Admin" {
		reports, err := h.service.GetAllReports()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Reports = reports

		requests, err := h.service.GetAllRequests()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.ModeratorRequests = requests

		users, err := h.service.GetAllUsers()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Users = users

		data.Form = models.UserLoginForm{}
		categories, err := h.service.GetCategories()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Categories = categories
	}

	notifications, err := h.service.GetUnreadNotifications(user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data.Notifications = notifications

	h.Render(w, http.StatusOK, "account.tmpl", data)
}

func (h *Handler) userRequestModerator(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
	}
	data.User = user

	if user.Role != "User" {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		return
	}

	err = h.service.RequestModeratorRole(user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) promoteUserToModerator(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	err = h.service.PromoteUserToModerator(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) denyModeratorRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	err = h.service.DenyModeratorRequest(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) demoteModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	userID, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil || userID < 1 || !UserExists(userID, users) {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	admin, err := h.service.GetUser(r)
	if err != nil || admin.Role != "Admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = h.service.DemoteModerator(userID)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
