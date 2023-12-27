package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskUserInfo() {
	defaultOutPath := getDefaultOutPath()

	fmt.Printf("Default Output folder: %s\nOutput folder (Just hit enter to use default path): ", defaultOutPath)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		panic(err)
	}

	input = strings.TrimSuffix(input, "\n")
	println(input)
	println(defaultOutPath)
}

func getDefaultOutPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir
}
