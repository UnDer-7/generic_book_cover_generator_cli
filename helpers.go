package main

import "os"

func (app *AppContext) currentPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		app.logger.Warn("Error getting current path")
		return "", err
	}

	return path, nil
}
