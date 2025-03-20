package main

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type NumberWindow struct {
	mu           sync.Mutex
	numbers      []int
	isFull       bool
	currentIndex int
}

type Response struct {
	WindowPrevState []int   `json:"windowPrevState"`
	WindowCurrState []int   `json:"windowCurrState"`
	Numbers         []int   `json:"numbers"`
	Avg             float64 `json:"avg"`
}

func (app *application) handleNumbers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	numberID := ps.ByName("numberId")

	if numberID != "p" && numberID != "f" && numberID != "e" && numberID != "r" {
		http.Error(w, "Invalid number ID. Valid IDs are 'p' (prime), 'f' (fibonacci), 'e' (even), or 'r' (random)", http.StatusBadRequest)
		return
	}

	prevState := app.window.GetCurrentState()

	numbers, err := app.fetchNumbers(numberID)
	if err != nil {
		app.logger.Printf("Error fetching numbers: %v", err)
		app.sendResponse(w, prevState, app.window.GetCurrentState(), []int{}, app.window.CalculateAverage())
		return
	}
	for _, num := range numbers {
		app.window.AddNumber(num)
	}
	app.sendResponse(w, prevState, app.window.GetCurrentState(), numbers, app.window.CalculateAverage())
}

func NewNumberWindow(size int) *NumberWindow {
	return &NumberWindow{
		numbers:      make([]int, size),
		currentIndex: 0,
		isFull:       false,
	}
}

func (w *NumberWindow) AddNumber(num int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := 0; i < len(w.numbers); i++ {
		if w.isFull || i < w.currentIndex {
			if w.numbers[i] == num {
				return false
			}
		}
	}

	w.numbers[w.currentIndex] = num
	w.currentIndex = (w.currentIndex + 1) % len(w.numbers)

	if w.currentIndex == 0 {
		w.isFull = true
	}

	return true
}

func (w *NumberWindow) GetCurrentState() []int {
	w.mu.Lock()
	defer w.mu.Unlock()

	var result []int

	if w.isFull {
		result = make([]int, len(w.numbers))
		for i := 0; i < len(w.numbers); i++ {
			index := (w.currentIndex + i) % len(w.numbers)
			result[i] = w.numbers[index]
		}
	} else {
		result = make([]int, w.currentIndex)
		copy(result, w.numbers[:w.currentIndex])
	}

	return result
}

func (w *NumberWindow) CalculateAverage() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentIndex == 0 && !w.isFull {
		return 0.0
	}

	sum := 0
	count := 0

	if w.isFull {
		count = len(w.numbers)
		for _, num := range w.numbers {
			sum += num
		}
	} else {
		count = w.currentIndex
		for i := 0; i < w.currentIndex; i++ {
			sum += w.numbers[i]
		}
	}

	return float64(sum) / float64(count)
}
