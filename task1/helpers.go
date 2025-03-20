package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetterNumbers struct {
	Numbers []int `json:"numbers"`
}

func (app *application) fetchNumbers(numberType string) ([]int, error) {
	client := http.Client{
		Timeout: time.Duration(app.config.maxResponseMs) * time.Millisecond,
	}
	num := ""
	if numberType == "p" {
		num = "primes"
	} else if numberType == "e" {
		num = "even"
	} else if numberType == "f" {
		num = "fibo"
	} else {
		num = "rand"
	}

	url := fmt.Sprintf("http://20.244.56.144/test/%s", num)
	app.logger.Printf("Fetching numbers from: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQyNDc3NjEyLCJpYXQiOjE3NDI0NzczMTIsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImVmYTMxZGFhLTNhNTAtNGFiYS1iYTRkLWM4MjE1YTg3MWMyZSIsInN1YiI6InNnOTg5MUBzcm1pc3QuZWR1LmluIn0sImNvbXBhbnlOYW1lIjoiU1JNSVNUIiwiY2xpZW50SUQiOiJlZmEzMWRhYS0zYTUwLTRhYmEtYmE0ZC1jODIxNWE4NzFjMmUiLCJjbGllbnRTZWNyZXQiOiJETklsT3VWZlpaTXNIZnVpIiwib3duZXJOYW1lIjoiU2FoaWwgR3VwdGEiLCJvd25lckVtYWlsIjoic2c5ODkxQHNybWlzdC5lZHUuaW4iLCJyb2xsTm8iOiJSQTIyMTEwMDMwMTE4MDkifQ.FzOMrApDt6R8OJvgTK9-Sbw-844oCbe4to1BADQUbvw")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("third-party service returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var thirdPartyResp GetterNumbers
	if err := json.Unmarshal(body, &thirdPartyResp); err != nil {
		return nil, err
	}

	return thirdPartyResp.Numbers, nil
}

func (app *application) sendResponse(w http.ResponseWriter, prevState, currState, numbers []int, avg float64) {
	response := Response{
		WindowPrevState: prevState,
		WindowCurrState: currState,
		Numbers:         numbers,
		Avg:             avg,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		app.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
