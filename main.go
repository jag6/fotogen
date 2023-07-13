package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jag6/fotogen/controllers"
	"github.com/jag6/fotogen/models"
	"github.com/jag6/fotogen/templates"
	"github.com/jag6/fotogen/views"
)

func main() {
	r := chi.NewRouter()

	//static files
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// //test
	// r.Get("/test", controllers.StaticHandler(views.Must(views.Parse("templates/base.html", "templates/components/header.html", "test.html"))))

	//homepage
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/home.html"))))

	//contact page
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/contact.html"))))

	//faq page
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "base.html", "pages/faq.html"))))

	//postgres connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	//signup page
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "base.html", "users/signup.html"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	//signin page
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "base.html", "users/signin.html"))
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	//logout
	r.Post("/signout", usersC.SignOut)

	//cookies
	r.Get("/users/me", usersC.CurrentUser)

	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("server running on " + "http://localhost:3000")
	csrfKey := ""
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))
	http.ListenAndServe(":3000", csrfMw(r))
}
