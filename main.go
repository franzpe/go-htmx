package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Increase(amount int) {
	c.mu.Lock()
	c.value = c.value + amount
	c.mu.Unlock()
}

func (c *Counter) Decrease(amount int) {
	c.mu.Lock()
	c.value = c.value - amount
	c.mu.Unlock()
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type CtxData struct {
	template *template.Template
	counter  *Counter
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	appData := r.Context().Value("ctxData").(CtxData)

	data := map[string]int{
		"CounterValue": appData.counter.GetValue(),
	}

	appData.template.ExecuteTemplate(w, "index.html", data)
}

func HandleIncrease(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		amount = 1
	}

	appData := r.Context().Value("ctxData").(CtxData)

	appData.counter.Increase(amount)

	data := map[string]int{
		"CounterValue": appData.counter.GetValue(),
	}

	appData.template.ExecuteTemplate(w, "counter.html", data)
}

func HandleDecrease(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		amount = 1
	}

	appData := r.Context().Value("ctxData").(CtxData)

	appData.counter.Decrease(amount)

	data := map[string]int{
		"CounterValue": appData.counter.GetValue(),
	}

	appData.template.ExecuteTemplate(w, "counter.html", data)
}

func MainCtx(counter *Counter, tmpl *template.Template) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "ctxData", CtxData{counter: counter, template: tmpl})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	counter := &Counter{}
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	tmpl, err := template.ParseGlob("./public/views/*.html")
	if err != nil {
		log.Fatalf("unable to parse templates %e\n", err)
	}

	r.Use(MainCtx(counter, tmpl))

	r.Get("/", HandleGet)
	r.Post("/increase", HandleIncrease)
	r.Post("/decrease", HandleDecrease)

	port := 3000
	log.Printf("Server has been spawned at port %d", port)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), r)
}
