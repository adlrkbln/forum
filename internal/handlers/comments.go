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

	comment_id, err := strconv.Atoi(r.FormValue("CommentId"))
	if err != nil || comment_id <= 0 || !CommentExists(comment_id, comments) {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	author, err := h.service.GetCommentAuthor(comment_id)
	if err != nil {
		h.ServerError(w, err)
	}
	if user.Role != "Admin" && user.Id != author {
		h.ClientError(w, http.StatusForbidden)
		return
	}

	if err := h.service.DeleteComment(comment_id); err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (h *Handler) commentEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.NotFound(w)
			return
		}
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		comment, err := h.service.GetComment(id)
		if err != nil {
			h.ServerError(w, err)
		}
		data.Comment = comment
		data.Form = models.Comment{
			Content: comment.Content,
		}
		data.Form = models.PostCreateForm{}
		h.Render(w, http.StatusOK, "edit_comment.tmpl", data)
	case http.MethodPost:
		h.commentEditPost(w, r)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) commentEditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil || id < 1 {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.CommentCreateForm{
		Content: r.PostForm.Get("content"),
	}

	form.CheckField(validate.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "edit_comment.tmpl", data)
		return
	}
	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = h.service.UpdateComment(id, form, data)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	comment, err := h.service.GetComment(id)
	if err != nil {
		h.ServerError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", comment.PostId), http.StatusSeeOther)
}
