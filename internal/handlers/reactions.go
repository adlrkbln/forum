package handlers

import (
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
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
	err = h.service.AddLikePost(post_id, user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", post_id), http.StatusSeeOther)
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
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
	err = h.service.AddDislikePost(post_id, user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", post_id), http.StatusSeeOther)
}

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	comments, err := h.service.GetAllComments()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	comment_idStr := r.FormValue("CommentId")
	comment_id, err := strconv.Atoi(comment_idStr)
	if err != nil || comment_id < 1 || !CommentExists(comment_id, comments) {
		h.NotFound(w)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	err = h.service.AddLikeComment(comment_id, user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	comments, err := h.service.GetAllComments()
	if err != nil {
		h.ServerError(w, err)
		return
	}
	comment_idStr := r.FormValue("CommentId")
	comment_id, err := strconv.Atoi(comment_idStr)
	if err != nil || comment_id < 1 || !CommentExists(comment_id, comments) {
		h.NotFound(w)
		return
	}

	user, err := h.service.GetUser(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	err = h.service.AddDislikeComment(comment_id, user.Id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func PostExists(postId int, posts []*models.Post) bool {
	for _, post := range posts {
		if post.Id == postId {
			return true
		}
	}
	return false
}

func CommentExists(commentId int, comments []*models.Comment) bool {
	for _, comment := range comments {
		if comment.Id == commentId {
			return true
		}
	}
	return false
}

func UserExists(userId int, users []*models.User) bool {
	for _, user := range users {
		if user.Id == userId {
			return true
		}
	}
	return false
}

func NotificationExists(notificationId int, notifications []*models.Notification) bool {
	for _, notification := range notifications {
		if notification.Id == notificationId {
			return true
		}
	}
	return false
}
