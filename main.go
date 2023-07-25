package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jag6/fotogen/controllers"
	"github.com/jag6/fotogen/migrations"
	"github.com/jag6/fotogen/models"
	"github.com/jag6/fotogen/templates"
	"github.com/jag6/fotogen/views"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	//psql
	cfg.PSQL = models.DefaultPostgresConfig()

	//smtp
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	//csrf
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = false

	//server
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	//database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	//services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	galleryService := &models.GalleryService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	//middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	//controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "base.html", "users/sign-up.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "base.html", "users/sign-in.html"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "base.html", "users/forgot-pw.html"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "base.html", "users/check-email.html"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "base.html", "users/reset-pw.html"))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.Index = views.Must(views.ParseFS(templates.FS, "base.html", "galleries/index.html"))
	galleriesC.Templates.New = views.Must(views.ParseFS(templates.FS, "base.html", "galleries/new.html"))
	galleriesC.Templates.Show = views.Must(views.ParseFS(templates.FS, "base.html", "galleries/show.html"))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(templates.FS, "base.html", "galleries/edit.html"))

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

	//user form pages
	r.Group(func(r chi.Router) {
		r.Use(umw.RedirectIfSignedIn)
		r.Get("/sign-up", usersC.New)
		r.Get("/sign-in", usersC.SignIn)
		r.Get("/forgot-pw", usersC.ForgotPassword)
		r.Get("/reset-pw", usersC.ResetPassword)
	})

	//sign up page
	r.Post("/users", usersC.Create)

	//sign in page
	r.Post("/sign-in", usersC.ProcessSignIn)

	//logout
	r.Post("/sign-out", usersC.SignOut)

	//forgot pw page
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	//reset pw page
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	//user page
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	//gallery pages
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/media/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			//index page
			r.Get("/", galleriesC.Index)
			//new page
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			//edit page
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			//delete gallery, image
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/media/{filename}/delete", galleriesC.DeleteImage)
		})
	})

	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	//start server
	fmt.Println("server running on " + "http://localhost" + cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
