package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func (app *AppContext) init() {
	defer app.measureExecutionTime()()

	app.createOutputCoverDir()
	app.validateIfOutputDirExists()

	names := app.getFileNames()

	if len(app.customCoverName) > 0 {
		app.processSingleChapterCustom()
		return
	}

	if app.useMultiThreading {
		app.runInMultiThreadMode(names)
		return
	}

	app.runInSingleThreadMode(names)
}

func (app *AppContext) runInMultiThreadMode(fileNames []string) {
	for _, name := range fileNames {
		app.wg.Add(1)
		go func(fileName string) {
			defer app.wg.Done()
			app.processSingleChapter(fileName)
		}(name)
	}

	app.wg.Wait()
}

func (app *AppContext) runInSingleThreadMode(fileNames []string) {
	for _, name := range fileNames {
		// ToDo: Add Option for when a processing of file throws an error it doesnt stop the whole script
		app.processSingleChapter(name)
	}
}

func (app *AppContext) processSingleChapterCustom() {
	app.logger.Info(fmt.Sprintf("processing custom cover %s", app.customCoverName))

	// ToDo: fix bug. When passing: -custom="Last of the day \n now plz" the "\n" is ignored and doesnt break line
	app.generateImg(app.customCoverName, "custom_cover")

	app.logger.Info(fmt.Sprintf("finished processing custom cover %s", app.customCoverName))
}

func (app *AppContext) processSingleChapter(fileName string) {
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
	app.generateImg(imageText, fileNameWithoutExtension)

	app.logger.Info(fmt.Sprintf("finished processing chapter %s", chapterNumber))
}
