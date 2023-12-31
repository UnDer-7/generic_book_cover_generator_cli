package main

import (
	"context"
	"embed"
	"flag"
	"log/slog"
	"os"
	"sync"
)

type AppPaths struct {
	bookFolder       string
	bookCoversOutput string
	font             string
	backgroundImage  string
}

type AppContext struct {
	path                             *AppPaths
	resources                        embed.FS
	logger                           *slog.Logger
	regexChapterNumberPrefix         string
	useMultiThreading                bool
	shouldCreateBookCoversOutput     bool
	wg                               sync.WaitGroup
	context                          context.Context
	skipFileWhenChapterNumberNotFund bool
	processOnlyBooksWithExtension    string
	customCoverName                  string
}

//go:embed assets
var assets embed.FS

func main() {
	appCtx := &AppContext{
		path:      &AppPaths{},
		resources: assets,
		context:   context.Background(),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}

	err := configureFlags(appCtx)
	if err != nil {
		panic(err)
	}

	appCtx.logger.Info("Starting script")
	appCtx.init()
	appCtx.logger.Info("Successfully finished script")
}

func configureFlags(appCtx *AppContext) error {
	defaultPath, err := appCtx.currentPath()
	defaultChapterPrefix := "Chapter_"
	if err != nil {
		return err
	}

	// todo: config log level
	flag.StringVar(&appCtx.path.bookFolder, "book", defaultPath, "Path where the books are located")
	flag.StringVar(&appCtx.path.bookCoversOutput, "out", defaultPath, "Folder path where the generated covers will be")
	flag.StringVar(&appCtx.regexChapterNumberPrefix, "prefix", defaultChapterPrefix, "Name that prefix the chapter number")
	flag.BoolVar(&appCtx.useMultiThreading, "t", false, "Whether to create the chapters covers using multi threading or not. TRUE=Multi Thread | FALSE=Single Thread")
	flag.BoolVar(&appCtx.shouldCreateBookCoversOutput, "c", false, "Create output cover folder if does not exists")
	flag.BoolVar(&appCtx.skipFileWhenChapterNumberNotFund, "s", true, "Skip file when chapter number is not found")
	flag.StringVar(&appCtx.processOnlyBooksWithExtension, "e", ".cbz", "File extensions to be used in the processing (pass with dot, ex: .cbz, .jpg, etc...)")
	flag.StringVar(&appCtx.customCoverName, "custom", "", "A custom cover name to be used in the image generation")

	flag.Parse()

	appCtx.path.font = "assets/font/Merriweather-Black.ttf"
	appCtx.path.backgroundImage = "assets/background/black_background.jpg"

	return nil
}
