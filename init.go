package main

import (
	"path/filepath"
	"strings"
)

func (app *AppContext) init() {
	defer app.measureExecutionTime()()

	app.createOutputCoverDir()
	app.validateIfOutputDirExists()

	names := app.getFileNames()

	for _, name := range names {
		if app.useMultiThreading {
			app.runInMultiThreadMode(name)
		} else {
			app.processSingleChapter(name)
		}

	}

}

func (app *AppContext) getChapterNumberAndFileNameWithoutExtension(fileName string) (chapterNumber, fileNameWithoutExtension string) {
	chapterNumber, err := app.extractChapterNumberFromFile(fileName)
	if err != nil {
		panic(err)
	}

	extension := filepath.Ext(fileName)
	fileNameWithoutExtension = strings.TrimSuffix(fileName, extension)

	return chapterNumber, fileNameWithoutExtension
}

func (app *AppContext) processSingleChapter(fileName string) {
	chapterNumber, fileNameWithoutExtension := app.getChapterNumberAndFileNameWithoutExtension(fileName)

	app.generateImg(chapterNumber, fileNameWithoutExtension)
}

func (app *AppContext) runInMultiThreadMode(fileName string) {
	// ToDo: Write code
	panic("Not supported yet")
}
