package main

import (
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrChapterNumberNotFound = errors.New("could not find chapter number in the file name")
)

func (app *AppContext) openBackgroundImage() fs.File {
	file, err := app.resources.Open(app.path.backgroundImage)
	if err != nil {
		app.logger.Warn("Error while opening black_background")
		panic(err)
	}

	return file
}

func (app *AppContext) decodeToJpeg(file fs.File) image.Image {
	img, err := jpeg.Decode(file)
	if err != nil {
		app.logger.Warn("Error while decoding black_background")
		panic(err)
	}
	return img
}

func (app *AppContext) LoadFont() *truetype.Font {
	fontBytes, err := app.resources.ReadFile(app.path.font)
	if err != nil {
		app.logger.Warn("Error while reading font file")
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		app.logger.Warn("Error while parsing font file")
		panic(err)
	}

	return font
}

func (app *AppContext) getFileNames() []string {
	entries, err := os.ReadDir(app.path.bookFolder)
	if err != nil {
		app.logger.Warn("Error while reading book folder")
		panic(err)
	}

	var fileNames []string
	var covers []string

	for _, entry := range entries {
		fileName := entry.Name()
		fileExtension := filepath.Ext(fileName)
		switch fileExtension {
		case app.bookCoverOutputExtension:
			covers = append(covers, fileName)
		case app.processOnlyBooksWithExtension:
			fileNames = append(fileNames, fileName)
		}
	}

	var validFileNames []string

	canAdd := func(fileNameSanitized string) bool {
		for _, cover := range covers {
			coverSanitized := strings.TrimSuffix(cover, filepath.Ext(app.bookCoverOutputExtension))

			if coverSanitized == fileNameSanitized {
				return false
			}
		}
		return true
	}

	for _, fileName := range fileNames {
		fileNameSanitized := strings.TrimSuffix(fileName, filepath.Ext(app.processOnlyBooksWithExtension))
		if canAdd(fileNameSanitized) {
			validFileNames = append(validFileNames, fileName)
		}
	}

	return validFileNames
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
