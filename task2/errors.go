package main

import (
	"net/http"
)

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}
