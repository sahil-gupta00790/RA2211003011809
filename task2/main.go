package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	testServerBaseURL = "http://20.244.56.144/test"
	cacheTTL          = 60 * time.Minute
)

type Cache struct {
	mu          sync.RWMutex
	users       map[string]string
	posts       map[string][]Post
	comments    map[string][]Comment
	topUsers    []UserWithPostCount
	topPosts    []PostWithCommentCount
	latestPosts []Post
	lastUpdated map[string]time.Time
}
type Post struct {
	ID      string `json:"id"`
	UserID  string `json:"userid"`
	Content string `json:"content"`
}

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"postid"`
	Content string `json:"content"`
}

type UserWithPostCount struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"postCount"`
}

type PostWithCommentCount struct {
	ID           string `json:"id"`
	UserID       string `json:"userid"`
	Content      string `json:"content"`
	CommentCount int    `json:"commentCount"`
}

var cache = Cache{
	users:       make(map[string]string),
	posts:       make(map[string][]Post),
	comments:    make(map[string][]Comment),
	lastUpdated: make(map[string]time.Time),
}

func (c *Cache) isCacheValid(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	lastUpdated, exists := c.lastUpdated[key]
	if !exists {
		return false
	}
	return time.Since(lastUpdated) < cacheTTL
}

func main() {
	r := httprouter.New()

	r.GET("/api/users/top", topUsersHandler)
	r.GET("/api/posts", postsHandler)
	r.NotFound = http.HandlerFunc(notFoundResponse)

	// Start server
	port := "8745"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
