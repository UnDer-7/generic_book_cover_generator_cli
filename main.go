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
	path                         *AppPaths
	resources                    embed.FS
	logger                       *slog.Logger
	regexChapterNumberPrefix     string
	useMultiThreading            bool
	shouldCreateBookCoversOutput bool
}

//go:embed assets
var assets embed.FS

func main() {
	appCtx := &AppContext{
		path:      &AppPaths{},
		resources: assets,
		logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	err := configureFlags(appCtx)
	if err != nil {
		panic(err)
	}

	appCtx.init()
}

func configureFlags(appCtx *AppContext) error {
	defaultPath, err := appCtx.currentPath()
	defaultChapterPrefix := "Chapter_"
	if err != nil {
		return err
	}

	flag.StringVar(&appCtx.path.bookFolder, "book", defaultPath, "Path where the books are located")
	flag.StringVar(&appCtx.path.bookCoversOutput, "out", defaultPath, "Folder path where the generated covers will be")
	flag.StringVar(&appCtx.regexChapterNumberPrefix, "prefix", defaultChapterPrefix, "Name that prefix the chapter number")
	flag.BoolVar(&appCtx.useMultiThreading, "t", false, "Whether to create the chapters covers using multi threading or not. TRUE=Multi Thread | FALSE=Single Thread")
	flag.BoolVar(&appCtx.shouldCreateBookCoversOutput, "c", false, "Create output cover folder if does not exists")

	flag.Parse()

	appCtx.path.font = "assets/font/Merriweather-Black.ttf"
	appCtx.path.backgroundImage = "assets/background/black_background.jpg"

	return nil
}
