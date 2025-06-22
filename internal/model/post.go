package model

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type NewPost struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}
