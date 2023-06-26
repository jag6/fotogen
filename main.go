package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/jag6/fotogen/controllers"
	"github.com/jag6/fotogen/views"
)

func main() {
	r := chi.NewRouter()

	//homepage
	r.Get("/", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "home.html")))))

	//contact page
	r.Get("/contact", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "contact.html")))))

	//faq page
	r.Get("/faq", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "faq.html")))))

	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("server running on " + "http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
