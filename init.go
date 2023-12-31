package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (app *AppContext) init() {
	defer app.measureExecutionTime()()

	app.createOutputCoverDir()
	app.validateIfOutputDirExists()

	names := app.getFileNames()

	if app.useMultiThreading {
		app.runInMultiThreadMode(names)
	} else {
		app.runInSingleThreadMode(names)
	}

	app.wg.Wait()
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
		app.processSingleChapter(name)
	}
}

func (app *AppContext) processSingleChapter(fileName string) {
	chapterNumber, err := app.extractChapterNumberFromFile(fileName)
	if err != nil {
		panic(err)
	}

	app.logger.Info(fmt.Sprintf("processing chapter %s", chapterNumber))
	extension := filepath.Ext(fileName)
	fileNameWithoutExtension := strings.TrimSuffix(fileName, extension)

	app.generateImg(chapterNumber, fileNameWithoutExtension)
	app.logger.Info(fmt.Sprintf("finished processing chapter %s", chapterNumber))
}
