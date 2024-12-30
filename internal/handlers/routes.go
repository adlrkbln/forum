package handlers

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/liked", h.RequireAuthentication(h.likedPosts))
	mux.HandleFunc("/post/create", h.RequireAuthentication(h.postCreate))
	mux.HandleFunc("/post/comment", h.RequireAuthentication(h.commentPost))
	mux.HandleFunc("/post/like", h.RequireAuthentication(h.likePost))
	mux.HandleFunc("/post/dislike", h.RequireAuthentication(h.dislikePost))

	mux.HandleFunc("/comment/like", h.RequireAuthentication(h.likeComment))
	mux.HandleFunc("/comment/dislike", h.RequireAuthentication(h.dislikeComment))

	mux.HandleFunc("/user/signup", h.userSignup)
	mux.HandleFunc("/user/login", h.userLogin)
	mux.HandleFunc("/user/logout", h.RequireAuthentication(h.userLogoutPost))
	mux.HandleFunc("/user/profile", h.RequireAuthentication(h.accountPageGet))

	mux.HandleFunc("/moderator/request", h.RequireAuthentication(h.userRequestModerator))
	mux.HandleFunc("/moderator/report", h.RequireModerator(h.reportPost))
	mux.HandleFunc("/admin/delete-post", h.RequireAdmin(h.deletePost))
	mux.HandleFunc("/admin/ignore-report", h.RequireAdmin(h.ignoreReport))
	mux.HandleFunc("/admin/promote", h.RequireAdmin(h.promoteUserToModerator))
	mux.HandleFunc("/admin/deny", h.RequireAdmin(h.denyModeratorRequest))

	return h.RecoverPanic(h.LogRequest(SecureHeaders(mux)))
}
