package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/MirzaHilmi/JariyahMu/internal/database"
	"github.com/MirzaHilmi/JariyahMu/internal/env"
	"github.com/MirzaHilmi/JariyahMu/internal/smtp"
	"github.com/MirzaHilmi/JariyahMu/internal/version"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	app, err := bootstrap(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

	err = app.serveHTTP()
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	jwt struct {
		secretKey string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
}

type application struct {
	config config
	db     *database.DB
	logger *slog.Logger
	mailer *smtp.Mailer
	wg     sync.WaitGroup
}

func bootstrap(logger *slog.Logger) (*application, error) {
	var cfg config

	err := env.LoadEnv("JariyahMu")
	if err != nil {
		return nil, err
	}

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.httpPort = env.GetInt("HTTP_PORT", 8080)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "")
	cfg.db.dsn = env.GetString("DB_DSN", "user:pass@tcp(localhost:3306)/example?parseTime=true")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", false)
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "")
	cfg.smtp.host = env.GetString("SMTP_HOST", "")
	cfg.smtp.port = env.GetInt("SMTP_PORT", 0)
	cfg.smtp.username = env.GetString("SMTP_USERNAME", "")
	cfg.smtp.password = env.GetString("SMTP_PASSWORD", "")
	cfg.smtp.from = env.GetString("SMTP_FROM", "")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil, nil
	}

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	if err != nil {
		return nil, err
	}

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		mailer: mailer,
	}

	return app, nil
}
