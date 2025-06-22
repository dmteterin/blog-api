package storage

import (
	"blog-api/internal/model"
	"testing"

	"github.com/rs/zerolog"
)

const shouldSeed = false

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    model.NewPost
		expected model.Post
	}{
		{
			name:  "basic create",
			input: model.NewPost{Title: "Post 1", Content: "Content 1", Author: "Alice"},
			expected: model.Post{
				ID:      1,
				Title:   "Post 1",
				Content: "Content 1",
				Author:  "Alice",
			},
		},
		{
			name:  "second post",
			input: model.NewPost{Title: "Post 2", Content: "Content 2", Author: "Bob"},
			expected: model.Post{
				ID:      2,
				Title:   "Post 2",
				Content: "Content 2",
				Author:  "Bob",
			},
		},
	}

	logger := zerolog.Nop()
	store := NewInMemoryStorage(logger, shouldSeed)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := store.Create(tt.input)

			if post.ID != tt.expected.ID {
				t.Errorf("Expected ID %d, got %d", tt.expected.ID, post.ID)
			}
			if post.Title != tt.expected.Title {
				t.Errorf("Expected title %q, got %q", tt.expected.Title, post.Title)
			}
			if post.Content != tt.expected.Content {
				t.Errorf("Expected content %q, got %q", tt.expected.Content, post.Content)
			}
			if post.Author != tt.expected.Author {
				t.Errorf("Expected author %q, got %q", tt.expected.Author, post.Author)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	logger := zerolog.Nop()
	store := NewInMemoryStorage(logger, shouldSeed)

	postsToCreate := []model.NewPost{
		{Title: "One", Content: "1", Author: "A"},
		{Title: "Two", Content: "2", Author: "B"},
		{Title: "Three", Content: "3", Author: "C"},
	}

	for _, np := range postsToCreate {
		store.Create(np)
	}

	got := store.GetAll()
	if len(got) != len(postsToCreate) {
		t.Fatalf("Expected %d posts, got %d", len(postsToCreate), len(got))
	}

	for i, p := range got {
		expected := postsToCreate[len(postsToCreate)-1-i]
		if p.Title != expected.Title {
			t.Errorf("Expected title %q, got %q", expected.Title, p.Title)
		}
		if p.Content != expected.Content {
			t.Errorf("Expected content %q, got %q", expected.Content, p.Content)
		}
		if p.Author != expected.Author {
			t.Errorf("Expected author %q, got %q", expected.Author, p.Author)
		}
	}
}

func TestGetByID(t *testing.T) {
	logger := zerolog.Nop()
	store := NewInMemoryStorage(logger, shouldSeed)

	post := store.Create(model.NewPost{Title: "Title", Content: "Content", Author: "Author"})

	tests := []struct {
		name        string
		id          int
		expectedErr bool
		expectedID  int
	}{
		{"found", post.ID, false, post.ID},
		{"not found", 999, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetByID(tt.id)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error, got: %v", err)
				}
				if got.ID != tt.expectedID {
					t.Errorf("Expected ID %d, got %d", tt.expectedID, got.ID)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	logger := zerolog.Nop()
	store := NewInMemoryStorage(logger, shouldSeed)

	existing := store.Create(model.NewPost{Title: "Old", Content: "Old", Author: "Author"})

	tests := []struct {
		name        string
		id          int
		update      model.NewPost
		expectedErr bool
		expected    string
	}{
		{
			name:        "update success",
			id:          existing.ID,
			update:      model.NewPost{Title: "New", Content: "Updated", Author: "Author"},
			expectedErr: false,
			expected:    "New",
		},
		{
			name:        "not found",
			id:          999,
			update:      model.NewPost{Title: "None"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := store.Update(tt.id, tt.update)
			if tt.expectedErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if updated.Title != tt.expected {
					t.Errorf("Expected title %s, got %s", tt.expected, updated.Title)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	logger := zerolog.Nop()
	store := NewInMemoryStorage(logger, shouldSeed)

	post := store.Create(model.NewPost{Title: "ToDelete", Content: "Delete", Author: "Author"})

	tests := []struct {
		name        string
		id          int
		expectedErr bool
	}{
		{"delete success", post.ID, false},
		{"delete fail", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.id)
			if tt.expectedErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
