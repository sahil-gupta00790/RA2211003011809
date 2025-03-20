package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func topUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if cache.isCacheValid("topUsers") {
		cache.mu.RLock()
		topUsers := cache.topUsers
		cache.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": topUsers,
		})
		return
	}

	if !cache.isCacheValid("users") {
		if err := fetchUser(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to fetch users",
			})
			return
		}
	}
	postCounts := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	cache.mu.RLock()
	users := cache.users
	cache.mu.RUnlock()

	for userID := range users {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			var posts []Post
			var err error

			cacheKey := fmt.Sprintf("posts_%s", id)
			if cache.isCacheValid(cacheKey) {
				cache.mu.RLock()
				posts = cache.posts[id]
				cache.mu.RUnlock()
			} else {
				posts, err = fetchUP(id)
				if err != nil {
					log.Printf("Error fetching posts for user %s: %v", id, err)
					return
				}
			}

			mu.Lock()
			postCounts[id] = len(posts)
			mu.Unlock()
		}(userID)
	}

	wg.Wait()

	userIDs := make([]string, 0, len(postCounts))
	for userID := range postCounts {
		userIDs = append(userIDs, userID)
	}

	sort.Slice(userIDs, func(i, j int) bool {
		return postCounts[userIDs[i]] > postCounts[userIDs[j]]
	})

	topN := 5
	if len(userIDs) < topN {
		topN = len(userIDs)
	}
	topUsers := make([]UserWithPostCount, 0, topN)
	for i := 0; i < topN; i++ {
		userID := userIDs[i]
		topUsers = append(topUsers, UserWithPostCount{
			ID:        userID,
			Name:      users[userID],
			PostCount: postCounts[userID],
		})
	}

	cache.mu.Lock()
	cache.topUsers = topUsers
	cache.lastUpdated["topUsers"] = time.Now()
	cache.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": topUsers,
	})
}

func postsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	postType := queryValues.Get("type")
	if postType == "" {
		postType = "popular"
	}

	if postType != "popular" && postType != "latest" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid type parameter. Use 'popular' or 'latest'.",
		})
		return
	}
	cacheKey := fmt.Sprintf("%sPosts", postType)
	if cache.isCacheValid(cacheKey) {
		var posts interface{}
		cache.mu.RLock()
		if postType == "popular" {
			posts = cache.topPosts
		} else {
			posts = cache.latestPosts
		}
		cache.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"posts": posts,
		})
		return
	}
	if !cache.isCacheValid("users") {
		if err := fetchUser(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to fetch users",
			})
			return
		}
	}
	var allPosts []Post
	var wg sync.WaitGroup
	var mu sync.Mutex

	cache.mu.RLock()
	users := cache.users
	cache.mu.RUnlock()

	for userID := range users {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			var posts []Post
			var err error

			cacheKey := fmt.Sprintf("posts_%s", id)
			if cache.isCacheValid(cacheKey) {
				cache.mu.RLock()
				posts = cache.posts[id]
				cache.mu.RUnlock()
			} else {
				posts, err = fetchUP(id)
				if err != nil {
					log.Printf("Error fetching posts for user %s: %v", id, err)
					return
				}
			}

			mu.Lock()
			allPosts = append(allPosts, posts...)
			mu.Unlock()
		}(userID)
	}

	wg.Wait()

	if postType == "popular" {
		postCommentCounts := make(map[string]int)
		for _, post := range allPosts {
			wg.Add(1)
			go func(p Post) {
				defer wg.Done()

				var comments []Comment
				var err error

				cacheKey := fmt.Sprintf("comments_%s", p.ID)
				if cache.isCacheValid(cacheKey) {
					cache.mu.RLock()
					comments = cache.comments[p.ID]
					cache.mu.RUnlock()
				} else {
					comments, err = fetchPC(p.ID)
					if err != nil {
						log.Printf("Error fetching comments for post %s: %v", p.ID, err)
						return
					}
				}

				mu.Lock()
				postCommentCounts[p.ID] = len(comments)
				mu.Unlock()
			}(post)
		}

		wg.Wait()
		var maxCommentCount int
		for _, count := range postCommentCounts {
			if count > maxCommentCount {
				maxCommentCount = count
			}
		}

		var topPosts []PostWithCommentCount
		for _, post := range allPosts {
			if postCommentCounts[post.ID] == maxCommentCount {
				topPosts = append(topPosts, PostWithCommentCount{
					ID:           post.ID,
					UserID:       post.UserID,
					Content:      post.Content,
					CommentCount: maxCommentCount,
				})
			}
		}

		cache.mu.Lock()
		cache.topPosts = topPosts
		cache.lastUpdated["popularPosts"] = time.Now()
		cache.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"posts": topPosts,
		})
	} else {
		sort.Slice(allPosts, func(i, j int) bool {
			return allPosts[i].ID > allPosts[j].ID
		})
		topN := 5
		if len(allPosts) < topN {
			topN = len(allPosts)
		}
		latestPosts := allPosts[:topN]

		// Update cache
		cache.mu.Lock()
		cache.latestPosts = latestPosts
		cache.lastUpdated["latestPosts"] = time.Now()
		cache.mu.Unlock()

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"posts": latestPosts,
		})
	}
}
