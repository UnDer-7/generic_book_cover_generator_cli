package main

import (
	"errors"
	"fmt"
	"github.com/golang/freetype/truetype"
	"image"
	"path/filepath"
	"strings"
)

func (app *AppContext) init() {
	defer app.measureExecutionTime()()

	app.createOutputCoverDir()
	app.validateIfOutputDirExists()

	names := app.getFileNames()
	bgFile := app.openBackgroundImage()
	defer bgFile.Close()
	bgImage := app.decodeToJpeg(bgFile)
	font := app.LoadFont()

	if len(app.customCoverName) > 0 {
		app.processChapterCustom(bgImage, font)
		return
	}

	if app.useMultiThreading {
		app.runInMultiThreadMode(names, bgImage, font)
		return
	}

	app.runInSingleThreadMode(names, bgImage, font)
}

func (app *AppContext) runInMultiThreadMode(fileNames []string, bgImage image.Image, font *truetype.Font) {
	autoFormGoroutinesOpened := 0
	app.logger.Debug("starting multi threaded processing")
	for _, name := range fileNames {
		app.wg.Add(1)
		autoFormGoroutinesOpened = autoFormGoroutinesOpened + 1
		go func(fileName string, bgImageInner image.Image) {
			defer app.wg.Done()
			app.processChapter(fileName, bgImageInner, font)
		}(name, bgImage)
	}

	app.logger.Debug(fmt.Sprintf("Amout of Go Routines opened %d", autoFormGoroutinesOpened))
	app.wg.Wait()
	app.logger.Debug("finished multi threaded processing")
}

func (app *AppContext) runInSingleThreadMode(fileNames []string, bgImage image.Image, font *truetype.Font) {
	for _, name := range fileNames {
		// ToDo: Add Option for when a processing of file throws an error it doesnt stop the whole script
		app.processChapter(name, bgImage, font)
	}
}

func (app *AppContext) processChapterCustom(bgImage image.Image, font *truetype.Font) {
	app.logger.Info(fmt.Sprintf("processing custom cover %s", app.customCoverName))

	// ToDo: fix bug. When passing: -custom="Last of the day \n now plz" the "\n" is ignored and doesnt break line
	app.generateImg(app.customCoverName, "custom_cover", bgImage, font)

	app.logger.Info(fmt.Sprintf("finished processing custom cover %s", app.customCoverName))
}

func (app *AppContext) processChapter(fileName string, bgImage image.Image, font *truetype.Font) {
	chapterNumber, err := app.extractChapterNumberFromFile(fileName)
	if err != nil {
		if errors.Is(err, ErrChapterNumberNotFound) {
			if app.skipFileWhenChapterNumberNotFund {
				app.logger.Debug(fmt.Sprintf("could not extract chapter number from file name %s", fileName))
				return
			}
			panic(err)
		}
	}

	app.logger.Info(fmt.Sprintf("processing chapter %s", chapterNumber))
	extension := filepath.Ext(fileName)
	fileNameWithoutExtension := strings.TrimSuffix(fileName, extension)

	imageText := "Chapter\n\n" + chapterNumber
	app.generateImg(imageText, fileNameWithoutExtension, bgImage, font)

	app.logger.Info(fmt.Sprintf("finished processing chapter %s", chapterNumber))
}
