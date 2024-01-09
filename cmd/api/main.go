package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/MirzaHilmi/JariyahMu/internal/bootstrap"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	app, err := bootstrap.Bootstrap(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

	err = app.ServeHTTP()
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}
