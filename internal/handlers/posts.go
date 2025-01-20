package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/validate"
	"net/http"
	"strconv"
)

func (h *Handler) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.NotFound(w)
		return
	}

	post, err := h.service.GetPost(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFound(w)
		} else {
			h.ServerError(w, err)
		}
		return
	}
	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data.Post = post

	h.Render(w, http.StatusOK, "view.tmpl", data)
}

func (h *Handler) postCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		categories, err := h.service.GetCategories()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = models.PostCreateForm{}
		data.Categories = categories
		h.Render(w, http.StatusOK, "create.tmpl", data)
	case http.MethodPost:
		h.postCreatePost(w, r)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) postCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	categoryIdsStr := r.Form["categoryIds[]"]
	var categoryIds []int
	for _, idStr := range categoryIdsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		categories, err := h.service.GetCategories()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		for _, category := range categories {
			if category.Id == id {
				categoryIds = append(categoryIds, id)
				break
			}
		}
	}

	form := models.PostCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		CategoryIds: categoryIds,
	}

	form.CheckField(validate.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validate.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validate.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validate.CheckCategory(form.CategoryIds), "categoryIds[]", "Choose existing categories")

	if !form.Valid() {
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		categories, err := h.service.GetCategories()
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = form
		data.Categories = categories
		h.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	id, err := h.service.InsertPost(form, data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = h.service.PostCategoryPost(id, categoryIds)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}

func (h *Handler) reportPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	posts, err := h.service.GetAllPosts()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	post_idStr := r.FormValue("post_id")
	post_id, err := strconv.Atoi(post_idStr)
	if err != nil || post_id < 1 || !PostExists(post_id, posts) {
		h.NotFound(w)
		return
	}
	reason := r.FormValue("reason")

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	if user.Role != "Moderator" {
		h.ClientError(w, http.StatusForbidden)
		return
	}

	err = h.service.ReportPost(user.Id, post_id, reason)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	author, err := h.service.GetPostAuthor(post_id)
	if err != nil {
		h.ServerError(w, err)
	}

	if user.Role != "Admin" || user.Id == author {
		h.ClientError(w, http.StatusForbidden)
		return
	}
	err = h.service.DeletePost(post_id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) ignoreReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	report_idStr := r.FormValue("report_id")
	report_id, err := strconv.Atoi(report_idStr)
	if err != nil || report_id < 1 {
		h.NotFound(w)
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
	err = h.service.IgnoreReport(report_id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (h *Handler) postEdit(w http.ResponseWriter, r *http.Request) {
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
		post, err := h.service.GetPost(id)
		if err != nil {
			h.ServerError(w, err)
		}
		data.Post = post
		data.Form = models.PostCreateForm{
			Title:   post.Title,
			Content: post.Content,
		}
		data.Form = models.PostCreateForm{}
		h.Render(w, http.StatusOK, "edit_post.tmpl", data)
	case http.MethodPost:
		h.postEditPost(w, r)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) postEditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil || id < 1 {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.PostCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	form.CheckField(validate.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validate.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validate.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data, err := h.NewTemplateData(w, r)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "edit_post.tmpl", data)
		return
	}
	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = h.service.UpdatePost(id, form, data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}
