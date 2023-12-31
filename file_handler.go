package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrChapterNumberNotFound = errors.New("could not find chapter number in the file name")
)

func (app *AppContext) getFileNames() []string {
	entries, err := os.ReadDir(app.path.bookFolder)
	if err != nil {
		app.logger.Warn("Error while reading book folder")
		panic(err)
	}

	var fileNames []string

	for _, entry := range entries {
		fileName := entry.Name()
		fileExtension := filepath.Ext(fileName)
		if fileExtension == app.processOnlyBooksWithExtension {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames
}

func (app *AppContext) extractChapterNumberFromFile(fileName string) (string, error) {
	strRgx := app.regexChapterNumberPrefix + `(\d+)`
	pattern := regexp.MustCompile(strRgx)
	matches := pattern.FindStringSubmatch(fileName)

	if len(matches) > 1 {
		return matches[1], nil
	}

	app.logger.Warn(
		fmt.Sprintf("could not find chapter number in the file name [ regexPattnerUsed: %s | fileNameUsed: %s ]", strRgx, fileName))
	return "", ErrChapterNumberNotFound
}

func (app *AppContext) createOutputCoverDir() {
	if app.shouldCreateBookCoversOutput {
		exists, err := app.bookCoversOutputExists()
		if err != nil {
			panic("could no create book cover output, because an error occurred while checking if book covers output already exists")
		}

		if !exists {
			err := os.MkdirAll(app.path.bookCoversOutput, os.ModePerm)
			if err != nil {
				app.logger.Warn("Could not create book cover output")
				panic(err)
			}
		}
	}
}

func (app *AppContext) bookCoversOutputExists() (bool, error) {
	_, err := os.Stat(app.path.bookCoversOutput)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
