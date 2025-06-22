package storage

import (
	"blog-api/internal/model"
	"fmt"
	"sort"
	"sync"

	"github.com/rs/zerolog"
)

const postsPath = "./internal/storage/seed-posts/posts.json"

type inMemoryStorage struct {
	logger zerolog.Logger
	posts  map[int]model.Post
	nextID int
	mu     sync.Mutex
}

func NewInMemoryStorage(l zerolog.Logger, shouldSeed bool) *inMemoryStorage {
	s := &inMemoryStorage{
		logger: l,
		posts:  make(map[int]model.Post),
		nextID: 1,
	}

	if !shouldSeed {
		return s
	}

	posts, err := loadPostsFromJSON(postsPath)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to load posts from JSON")
		return s
	}

	for _, p := range posts {
		s.Create(p)
	}

	return s
}

func (s *inMemoryStorage) GetAll() []model.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	postList := make([]model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		postList = append(postList, post)
	}

	// Default sorting - by ID Desc
	sort.Slice(postList, func(i, j int) bool {
		return postList[i].ID > postList[j].ID
	})

	return postList
}

func (s *inMemoryStorage) GetByID(id int) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return model.Post{}, fmt.Errorf("post with ID %d not found in storage", id)
	}
	return post, nil
}

func (s *inMemoryStorage) Create(newPost model.NewPost) model.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := model.Post{
		ID:      s.nextID,
		Title:   newPost.Title,
		Content: newPost.Content,
		Author:  newPost.Author,
	}
	s.posts[post.ID] = post
	s.nextID++
	return post
}

func (s *inMemoryStorage) Update(id int, updatedPost model.NewPost) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return model.Post{}, fmt.Errorf("post with ID %d not found in storage", id)
	}

	post := model.Post{
		ID:      id,
		Title:   updatedPost.Title,
		Content: updatedPost.Content,
		Author:  updatedPost.Author,
	}
	s.posts[id] = post
	return post, nil
}

func (s *inMemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return fmt.Errorf("post with ID %d not found in storage", id)
	}

	delete(s.posts, id)
	return nil
}
