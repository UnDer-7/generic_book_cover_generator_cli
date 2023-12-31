package main

import (
	"fmt"
	"os"
	"time"
)

func (app *AppContext) currentPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		app.logger.Warn("Error getting current path")
		return "", err
	}

	return path, nil
}

func (app *AppContext) measureExecutionTime() func() {
	start := time.Now()
	app.logger.Info(fmt.Sprintf("Starting measuring execution time %v", start))
	return func() {
		endTime := time.Now()
		diffTIme := endTime.Sub(start)
		app.logger.Info(fmt.Sprintf(
			"Finished measuring execution time. Ended At: %v | Took: [ %v Minutes - %v Second - %v Millisecond - %v microsecond",
			endTime,
			diffTIme.Minutes(),
			diffTIme.Seconds(),
			diffTIme.Milliseconds(),
			diffTIme.Microseconds(),
		))
	}
}

func (app *AppContext) validateIfOutputDirExists() {
	exists, err := app.bookCoversOutputExists()
	if err != nil {
		app.logger.Warn("Could check if book covers output path exists")
	}

	if !exists {
		panic(fmt.Sprintf("book covers output path does not exists. | [ bookCoversOutput: %s ]", app.path.bookCoversOutput))
	}
}
