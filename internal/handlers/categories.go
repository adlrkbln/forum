package handlers

import (
	"forum/internal/models"
	"forum/internal/validate"
	"net/http"
	"strconv"
)

func (h *Handler) addCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var form models.CategoryCreateForm
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form.Name = r.PostForm.Get("category_name")
	form.CheckField(validate.NotBlank(form.Name), "name", "This field cannot be blank")

	if err := h.service.CreateCategory(form); err != nil || !form.Valid() {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))

	found := false
	categories, err := h.service.GetCategories()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	for _, category := range categories {
		if category.Id == categoryID {
			found = true
			break
		}
	}
	if err != nil || categoryID <= 0 || !found {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteCategory(categoryID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) markNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    err := r.ParseForm()
    if err != nil {
        h.ClientError(w, http.StatusBadRequest)
        return
    }
	notifications, err := h.service.GetNotifications()
	if err != nil {
		h.ServerError(w, err)
		return
	}

    notificationIdStr := r.PostForm.Get("notification_id")
    if notificationIdStr == "" {
        h.ClientError(w, http.StatusBadRequest)
        return
    }

	notificationId, err := strconv.Atoi(notificationIdStr)
	if err != nil || notificationId < 1 || !NotificationExists(notificationId, notifications) {
		h.NotFound(w)
		return
	}
    err = h.service.MarkNotificationAsRead(notificationId)
    if err != nil {
        h.ServerError(w, err)
        return
    }

    http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
