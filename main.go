package main

import (
	"embed"
	"flag"
	"log/slog"
	"os"
)

type AppPaths struct {
	bookFolder       string
	bookCoversOutput string
	font             string
	backgroundImage  string
}

type AppContext struct {
	path      *AppPaths
	resources embed.FS
	logger    *slog.Logger
}

//go:embed assets
var assets embed.FS

func main() {
	appCtx := &AppContext{
		path:      &AppPaths{},
		resources: assets,
		logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	err := configurePaths(appCtx)
	if err != nil {
		panic(err)
	}

	appCtx.init()
}

func configurePaths(appCtx *AppContext) error {
	defaultPath, err := appCtx.currentPath()
	if err != nil {
		return err
	}

	flag.StringVar(&appCtx.path.bookFolder, "book", defaultPath, "Path where the books are located")
	flag.StringVar(&appCtx.path.bookCoversOutput, "out", defaultPath, "Path where the generated covers will be")

	flag.Parse()

	appCtx.path.font = "assets/font/Merriweather-Black.ttf"
	appCtx.path.backgroundImage = "assets/background/black_background.jpg"

	return nil
}
