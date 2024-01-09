package tests

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/MirzaHilmi/JariyahMu/internal/bootstrap"
	"github.com/lmittmann/tint"
)

var app *bootstrap.Application

func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	var err error

	app, err = bootstrap.Bootstrap(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}
