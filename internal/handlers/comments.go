package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/validate"
	"net/http"
	"strconv"
)

func (h *Handler) commentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.GetAllPosts()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	post_idStr := r.FormValue("PostId")
	post_id, err := strconv.Atoi(post_idStr)
	if err != nil || post_id < 1 || !PostExists(post_id, posts) {
		h.NotFound(w)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
	}

	content := r.FormValue("content")

	form := models.CommentCreateForm{
		Content: r.PostForm.Get("content"),
	}

	form.CheckField(validate.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", post_id), http.StatusSeeOther)
		return
	}

	err = h.service.InsertComment(post_id, user.Id, content)
	if err != nil {
		h.ClientError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", post_id), http.StatusSeeOther)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
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
		h.ClientError(w, http.StatusForbidden)
		return
	}

	if err := h.service.DeleteComment(commentID); err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (h *Handler) commentEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Fetch comment details and render the edit form
	} else if r.Method == http.MethodPost {
		// Update the comment in the database
	}
}
