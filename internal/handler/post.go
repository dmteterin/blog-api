package handler

import (
	"blog-api/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

type PostService interface {
	GetPostsPaginated(limit, offset int) []model.Post
	GetPostByID(id int) (model.Post, error)
	CreatePost(newPost model.NewPost) model.Post
	UpdatePost(id int, updatedPost model.NewPost) (model.Post, error)
	DeletePost(id int) error
}

type PostHandler struct {
	service PostService
}

func NewPostHandler(s PostService) *PostHandler {
	return &PostHandler{
		service: s,
	}
}

func (h *PostHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit := defaultLimit
	offset := defaultOffset

	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	posts := h.service.GetPostsPaginated(limit, offset)
	respondWithJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, post)
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var newPost model.NewPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if newPost.Title == "" || newPost.Content == "" || newPost.Author == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields: title, content, author")
		return
	}

	createdPost := h.service.CreatePost(newPost)
	respondWithJSON(w, http.StatusCreated, createdPost)
}

func (h *PostHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var updatedPost model.NewPost
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	post, err := h.service.UpdatePost(id, updatedPost)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, post)
}

func (h *PostHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.service.DeletePost(id); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
