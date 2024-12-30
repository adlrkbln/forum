package handlers

import (
	"forum/internal/models"
	"forum/internal/validate"
	"net/http"
	"strconv"
)

func (h *Handler) AdminManageCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var form models.CategoryCreateForm
		err := r.ParseForm()
		if err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}

		form.Name = r.PostForm.Get("category_name")
		form.CheckField(validate.NotBlank(form.Name), "name", "This field cannot be blank")

		if err := h.service.CreateCategory(form); err != nil || !form.Valid() {
			data, err := h.NewTemplateData(w, r)
			if err != nil {
				h.ServerError(w, err)
				return
			}
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "admin.tmpl", data)
			return
		}

		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)

	case http.MethodDelete:
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

		if err != nil || categoryID <= 0 || !found{
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		if err := h.service.DeleteCategory(categoryID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
