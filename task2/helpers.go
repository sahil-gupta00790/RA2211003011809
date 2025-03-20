package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UsersResponse struct {
	Users map[string]string `json:"users"`
}

type PostsResponse struct {
	Posts []Post `json:"posts"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

func fetchUser() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", testServerBaseURL), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQyNDc3NjEyLCJpYXQiOjE3NDI0NzczMTIsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImVmYTMxZGFhLTNhNTAtNGFiYS1iYTRkLWM4MjE1YTg3MWMyZSIsInN1YiI6InNnOTg5MUBzcm1pc3QuZWR1LmluIn0sImNvbXBhbnlOYW1lIjoiU1JNSVNUIiwiY2xpZW50SUQiOiJlZmEzMWRhYS0zYTUwLTRhYmEtYmE0ZC1jODIxNWE4NzFjMmUiLCJjbGllbnRTZWNyZXQiOiJETklsT3VWZlpaTXNIZnVpIiwib3duZXJOYW1lIjoiU2FoaWwgR3VwdGEiLCJvd25lckVtYWlsIjoic2c5ODkxQHNybWlzdC5lZHUuaW4iLCJyb2xsTm8iOiJSQTIyMTEwMDMwMTE4MDkifQ.FzOMrApDt6R8OJvgTK9-Sbw-844oCbe4to1BADQUbvw")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var usersResp UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&usersResp); err != nil {
		return err
	}

	cache.mu.Lock()
	cache.users = usersResp.Users
	cache.lastUpdated["users"] = time.Now()
	cache.mu.Unlock()

	return nil
}

func fetchUP(userID string) ([]Post, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/posts", testServerBaseURL, userID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQyNDc3NjEyLCJpYXQiOjE3NDI0NzczMTIsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImVmYTMxZGFhLTNhNTAtNGFiYS1iYTRkLWM4MjE1YTg3MWMyZSIsInN1YiI6InNnOTg5MUBzcm1pc3QuZWR1LmluIn0sImNvbXBhbnlOYW1lIjoiU1JNSVNUIiwiY2xpZW50SUQiOiJlZmEzMWRhYS0zYTUwLTRhYmEtYmE0ZC1jODIxNWE4NzFjMmUiLCJjbGllbnRTZWNyZXQiOiJETklsT3VWZlpaTXNIZnVpIiwib3duZXJOYW1lIjoiU2FoaWwgR3VwdGEiLCJvd25lckVtYWlsIjoic2c5ODkxQHNybWlzdC5lZHUuaW4iLCJyb2xsTm8iOiJSQTIyMTEwMDMwMTE4MDkifQ.FzOMrApDt6R8OJvgTK9-Sbw-844oCbe4to1BADQUbvw")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var postsResp PostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&postsResp); err != nil {
		return nil, err
	}

	cache.mu.Lock()
	cache.posts[userID] = postsResp.Posts
	cache.lastUpdated[fmt.Sprintf("posts_%s", userID)] = time.Now()
	cache.mu.Unlock()

	return postsResp.Posts, nil
}

func fetchPC(postID string) ([]Comment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/posts/%s/comments", testServerBaseURL, postID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQyNDc3NjEyLCJpYXQiOjE3NDI0NzczMTIsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImVmYTMxZGFhLTNhNTAtNGFiYS1iYTRkLWM4MjE1YTg3MWMyZSIsInN1YiI6InNnOTg5MUBzcm1pc3QuZWR1LmluIn0sImNvbXBhbnlOYW1lIjoiU1JNSVNUIiwiY2xpZW50SUQiOiJlZmEzMWRhYS0zYTUwLTRhYmEtYmE0ZC1jODIxNWE4NzFjMmUiLCJjbGllbnRTZWNyZXQiOiJETklsT3VWZlpaTXNIZnVpIiwib3duZXJOYW1lIjoiU2FoaWwgR3VwdGEiLCJvd25lckVtYWlsIjoic2c5ODkxQHNybWlzdC5lZHUuaW4iLCJyb2xsTm8iOiJSQTIyMTEwMDMwMTE4MDkifQ.FzOMrApDt6R8OJvgTK9-Sbw-844oCbe4to1BADQUbvw")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commentsResp CommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentsResp); err != nil {
		return nil, err
	}

	cache.mu.Lock()
	cache.comments[postID] = commentsResp.Comments
	cache.lastUpdated[fmt.Sprintf("comments_%s", postID)] = time.Now()
	cache.mu.Unlock()

	return commentsResp.Comments, nil
}
