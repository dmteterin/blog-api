package handler

import (
	"blog-api/internal/handler/mocks"
	"blog-api/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestGetPosts(t *testing.T) {
	mockService := new(mocks.PostService)
	handler := NewPostHandler(mockService)

	expected := []model.Post{{ID: 1, Title: "Go", Content: "Post", Author: "Alice"}}
	mockService.On("GetAllPosts").Return(expected)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rr := httptest.NewRecorder()

	handler.getPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	mockService.AssertExpectations(t)
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		mockSetup    func(m *mocks.PostService)
		expectedCode int
	}{
		{
			name:    "valid id",
			idParam: "1",
			mockSetup: func(m *mocks.PostService) {
				m.On("GetPostByID", 1).Return(model.Post{ID: 1}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id format",
			idParam:      "abc",
			mockSetup:    func(m *mocks.PostService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "post not found",
			idParam: "99",
			mockSetup: func(m *mocks.PostService) {
				m.On("GetPostByID", 99).Return(model.Post{}, errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, "/posts/"+tt.idParam, nil)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			handler := NewPostHandler(mockService)
			router.HandleFunc("/posts/{id}", handler.getPost)

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, rr.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name         string
		body         model.NewPost
		expectStatus int
		mockSetup    func(m *mocks.PostService)
	}{
		{
			name: "valid payload",
			body: model.NewPost{
				Title:   "Title",
				Content: "Content",
				Author:  "Author",
			},
			expectStatus: http.StatusCreated,
			mockSetup: func(m *mocks.PostService) {
				m.On("CreatePost", mock.Anything).Return(model.Post{ID: 1})
			},
		},
		{
			name:         "missing fields",
			body:         model.NewPost{},
			expectStatus: http.StatusBadRequest,
			mockSetup:    func(m *mocks.PostService) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			tt.mockSetup(mockService)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler := NewPostHandler(mockService)
			handler.createPost(rr, req)

			if rr.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, rr.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		body         model.NewPost
		mockSetup    func(m *mocks.PostService)
		expectedCode int
	}{
		{
			name:    "valid update",
			idParam: "1",
			body: model.NewPost{
				Title:   "Updated",
				Content: "New content",
				Author:  "Author",
			},
			mockSetup: func(m *mocks.PostService) {
				m.On("UpdatePost", 1, mock.Anything).Return(model.Post{ID: 1}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id",
			idParam:      "abc",
			body:         model.NewPost{},
			mockSetup:    func(m *mocks.PostService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "post not found",
			idParam: "99",
			body:    model.NewPost{Title: "X", Content: "Y", Author: "Z"},
			mockSetup: func(m *mocks.PostService) {
				m.On("UpdatePost", 99, mock.Anything).Return(model.Post{}, errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			tt.mockSetup(mockService)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/posts/"+tt.idParam, bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			handler := NewPostHandler(mockService)
			router.HandleFunc("/posts/{id}", handler.updatePost)

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, rr.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		mockSetup    func(m *mocks.PostService)
		expectedCode int
	}{
		{
			name:    "valid delete",
			idParam: "1",
			mockSetup: func(m *mocks.PostService) {
				m.On("DeletePost", 1).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "invalid id",
			idParam:      "abc",
			mockSetup:    func(m *mocks.PostService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "post not found",
			idParam: "99",
			mockSetup: func(m *mocks.PostService) {
				m.On("DeletePost", 99).Return(errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.PostService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/posts/"+tt.idParam, nil)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			handler := NewPostHandler(mockService)
			router.HandleFunc("/posts/{id}", handler.deletePost)

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, rr.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}
