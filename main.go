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
	emailService := models.NewEmailService(cfg.SMTP)

	//middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

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

	//sign up page
	r.Get("/sign-up", usersC.New)
	r.Post("/users", usersC.Create)

	//sign in page
	r.Get("/sign-in", usersC.SignIn)
	r.Post("/sign-in", usersC.ProcessSignIn)

	//logout
	r.Post("/sign-out", usersC.SignOut)

	//forgot pw page
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	//reset pw page
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

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
	fmt.Println("server running on " + "http://localhost" + cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
