package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}
	categoryIdStr := r.URL.Query().Get("category")

	var categoryId int
	var err error

	if categoryIdStr != "" {
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	var posts []*models.Post
	if categoryId > 0 {
		posts, err = h.service.GetPostByCategory(categoryId)
	} else {
		posts, err = h.service.GetAllPosts()
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}

	categories, err := h.service.GetCategories()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data, err := h.NewTemplateData(w, r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	data.Posts = posts
	data.Categories = categories

	h.Render(w, http.StatusOK, "home.tmpl", data)
}
