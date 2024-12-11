package handlers

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/liked", h.requireAuthentication(h.likedPosts))
	mux.HandleFunc("/post/create", h.requireAuthentication(h.postCreate))
	mux.HandleFunc("/post/comment", h.requireAuthentication(h.commentPost))
	mux.HandleFunc("/post/like", h.requireAuthentication(h.likePost))
	mux.HandleFunc("/post/dislike", h.requireAuthentication(h.dislikePost))
	
	mux.HandleFunc("/comment/like", h.requireAuthentication(h.likeComment))
	mux.HandleFunc("/comment/dislike", h.requireAuthentication(h.dislikeComment))

	mux.HandleFunc("/user/signup", h.userSignup)
    mux.HandleFunc("/user/login", h.userLogin)
    mux.HandleFunc("/user/logout", h.requireAuthentication(h.userLogoutPost))

	
	return h.recoverPanic(h.logRequest(secureHeaders(mux)))
}
