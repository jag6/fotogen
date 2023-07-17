package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jag6/fotogen/controllers"
	"github.com/jag6/fotogen/migrations"
	"github.com/jag6/fotogen/models"
	"github.com/jag6/fotogen/templates"
	"github.com/jag6/fotogen/views"
)

func main() {
	//database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	//services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	//middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := ""
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	//controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "base.html", "users/signup.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "base.html", "users/signin.html"))

	//router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	//static files
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	//homepage
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/home.html"))))

	//contact page
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "base.html", "pages/contact.html"))))

	//faq page
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "base.html", "pages/faq.html"))))

	//signup page
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	//signin page
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	//logout
	r.Post("/signout", usersC.SignOut)

	//user page
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	//start server
	fmt.Println("server running on " + "http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
