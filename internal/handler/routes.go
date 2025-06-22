package handler

import "github.com/gorilla/mux"

func (h *PostHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/posts", h.getPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", h.getPost).Methods("GET")
	r.HandleFunc("/posts", h.createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", h.updatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", h.deletePost).Methods("DELETE")
}
