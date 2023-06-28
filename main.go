package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jag6/fotogen/controllers"
	"github.com/jag6/fotogen/static"
	"github.com/jag6/fotogen/templates"
	"github.com/jag6/fotogen/views"
)

func main() {
	r := chi.NewRouter()

	//static files
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(static.FS))))

	//homepage
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/home.html"))))

	//contact page
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/contact.html"))))

	//faq page
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "base.html", "pages/faq.html"))))

	//signup page
	usersC := controllers.Users{}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "base.html", "users/signup.html"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	//login page
	r.Get("/login", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "users/login.html"))))

	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("server running on " + "http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
