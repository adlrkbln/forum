package handlers

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/post/create", h.RequireAuthentication(h.postCreate))
	mux.HandleFunc("/post/comment", h.RequireAuthentication(h.commentPost))
	mux.HandleFunc("/post/like", h.RequireAuthentication(h.likePost))
	mux.HandleFunc("/post/dislike", h.RequireAuthentication(h.dislikePost))
	mux.HandleFunc("/post/edit", h.RequireAuthentication(h.postEdit))
	mux.HandleFunc("/post/delete", h.RequireAuthentication(h.deletePost))

	mux.HandleFunc("/comment/like", h.RequireAuthentication(h.likeComment))
	mux.HandleFunc("/comment/dislike", h.RequireAuthentication(h.dislikeComment))
	mux.HandleFunc("/comment/edit", h.RequireAuthentication(h.commentEdit))
	mux.HandleFunc("/comment/delete", h.RequireAuthentication(h.deleteComment))

	mux.HandleFunc("/user/signup", h.CheckGuest(h.userSignup))
	mux.HandleFunc("/user/login", h.CheckGuest(h.userLogin))
	mux.HandleFunc("/user/logout", h.RequireAuthentication(h.userLogoutPost))
	mux.HandleFunc("/user/profile", h.RequireAuthentication(h.accountPageGet))
	mux.HandleFunc("/notifications/read", h.RequireAuthentication(h.markNotificationRead))

	mux.HandleFunc("/moderator/request", h.RequireAuthentication(h.userRequestModerator))
	mux.HandleFunc("/moderator/report", h.RequireModerator(h.reportPost))

	mux.HandleFunc("/admin/ignore-report", h.RequireAdmin(h.ignoreReport))
	mux.HandleFunc("/admin/promote", h.RequireAdmin(h.promoteUserToModerator))
	mux.HandleFunc("/admin/deny", h.RequireAdmin(h.denyModeratorRequest))
	mux.HandleFunc("/admin/demote", h.RequireAdmin(h.demoteModerator))
	mux.HandleFunc("/admin/categories", h.RequireAdmin(h.addCategories))
	mux.HandleFunc("/admin/delete-category", h.RequireAdmin(h.deleteCategory))

	return h.RecoverPanic(h.LogRequest(SecureHeaders(mux)))
}
