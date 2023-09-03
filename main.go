package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Counter struct {
	value int        `json:"value"`
	mu    sync.Mutex `json:"mu"`
}

func (c *Counter) Increase() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Decrease() {
	c.mu.Lock()
	c.value--
	c.mu.Unlock()
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &Counter{}
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		tmpl, _ := template.ParseFiles("index.html")

		data := map[string]int{
			"CounterValue": counter.GetValue(),
		}

		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	r.Post("/increase", func(w http.ResponseWriter, _ *http.Request) {
		tmpl, _ := template.ParseFiles("counter.html")
		counter.Increase()
		data := map[string]int{
			"CounterValue": counter.GetValue(),
		}
		tmpl.ExecuteTemplate(w, "counter.html", data)
	})

	r.Post("/decrease", func(w http.ResponseWriter, _ *http.Request) {
		tmpl, _ := template.ParseFiles("counter.html")
		counter.Decrease()

		data := map[string]int{
			"CounterValue": counter.GetValue(),
		}

		tmpl.ExecuteTemplate(w, "counter.html", data)
	})

	http.ListenAndServe("localhost:3000", r)
}
