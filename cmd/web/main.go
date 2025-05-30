package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashparshp/bookings/internal/config"
	"github.com/ashparshp/bookings/internal/driver"
	"github.com/ashparshp/bookings/internal/handlers"
	"github.com/ashparshp/bookings/internal/helpers"
	"github.com/ashparshp/bookings/internal/models"
	"github.com/ashparshp/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	portNumber := getPort()
	fmt.Println("Server running on port", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// Register custom session data types
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// Read flags
	inProduction := flag.Bool("production", true, "Run in production mode")
	useCache := flag.Bool("cache", true, "Use template caching")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPassword := flag.String("dbpassword", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL setting (disable, prefer, require)")

	// Email configuration flags
    mailHost := flag.String("mailhost", "localhost", "SMTP host")
    mailPort := flag.Int("mailport", 1025, "SMTP port")
    mailUsername := flag.String("mailusername", "", "SMTP username")
    mailPassword := flag.String("mailpassword", "", "SMTP password")
    mailEncryption := flag.String("mailencryption", "none", "SMTP encryption (none, tls, ssl)")
    mailFromAddress := flag.String("mailfrom", "noreply@bookings.com", "Mail from address")
    mailFromName := flag.String("mailfromname", "Bookings", "Mail from name")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		log.Println("Database name and user must be provided")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Store email config in app
    app.MailConfig = config.MailConfig{
        Host:       *mailHost,
        Port:       *mailPort,
        Username:   *mailUsername,
        Password:   *mailPassword,
        Encryption: *mailEncryption,
        FromAddress: *mailFromAddress,
        FromName:   *mailFromName,
    }

	app.InProduction = *inProduction
	app.UseCahce = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Database connection
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL,
	)

	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache:", err)
		return nil, err
	}
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return ":" + port
}