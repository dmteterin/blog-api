package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func (h *PostHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/posts", h.getPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", h.getPost).Methods("GET")
	r.HandleFunc("/posts", h.createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", h.updatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", h.deletePost).Methods("DELETE")
}

func NewRouter(postHandler *PostHandler, logger zerolog.Logger) http.Handler {
	r := mux.NewRouter()
	postHandler.RegisterRoutes(r)

	loggingMiddleware := LoggingMiddleware(logger)
	return loggingMiddleware(r)
}
