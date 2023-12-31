package main

import (
	"os"
)

func (app *AppContext) getFileNames() []string {
	entries, err := os.ReadDir(app.path.bookFolder)
	if err != nil {
		app.logger.Warn("Error while reading book folder")
		panic(err)
	}

	var fileName []string

	for _, entry := range entries {
		fileName = append(fileName, entry.Name())
	}

	return fileName
}
