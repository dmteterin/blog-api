package storage

import (
	"blog-api/internal/model"
	"encoding/json"
	"os"
)

func loadPostsFromJSON(path string) ([]model.NewPost, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	type postsWrapper struct {
		Posts []struct {
			// Ignore ID
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
		} `json:"posts"`
	}

	var pw postsWrapper
	if err := json.Unmarshal(data, &pw); err != nil {
		return nil, err
	}

	newPosts := make([]model.NewPost, len(pw.Posts))
	for i, p := range pw.Posts {
		newPosts[i] = model.NewPost{
			Title:   p.Title,
			Content: p.Content,
			Author:  p.Author,
		}
	}
	return newPosts, nil
}
