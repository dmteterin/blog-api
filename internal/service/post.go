package service

import "blog-api/internal/model"

type PostStorage interface {
	GetPaginated(limit, offset int) []model.Post
	GetByID(id int) (model.Post, error)
	Create(newPost model.NewPost) model.Post
	Update(id int, updatedPost model.NewPost) (model.Post, error)
	Delete(id int) error
}

type postService struct {
	storage PostStorage
}

func NewPostService(s PostStorage) *postService {
	return &postService{
		storage: s,
	}
}

func (s *postService) GetPostsPaginated(limit, offset int) []model.Post {
	return s.storage.GetPaginated(limit, offset)
}

func (s *postService) GetPostByID(id int) (model.Post, error) {
	return s.storage.GetByID(id)
}

func (s *postService) CreatePost(newPost model.NewPost) model.Post {
	// This would be the place we could add business logic - validation, duplication checks, etc...
	return s.storage.Create(newPost)
}

func (s *postService) UpdatePost(id int, updatedPost model.NewPost) (model.Post, error) {
	return s.storage.Update(id, updatedPost)
}

func (s *postService) DeletePost(id int) error {
	return s.storage.Delete(id)
}
