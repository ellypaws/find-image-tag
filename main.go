package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nokusukun/roggy"
	"os"
	"strings"
)

var rogPrinter = roggy.Printer("main-service")

func main() {
	// Read all the filenames of image files in the image directory

	roggy.LogLevel = roggy.TypeDebug
	rogPrinter.Debug("Starting main service")

	data := DataSet{
		Images: make(map[string]Image),
	}

	data.promptOption()
}

func (data *DataSet) promptOption() {
	reader := bufio.NewReader(os.Stdin)

	roggy.Flush()
	fmt.Println("--- Image Captioning ---")
	fmt.Println("[A]dd files to the dataset")
	fmt.Println("[C]heck if captions exist")
	fmt.Println("[M]ove captions to the image files")
	fmt.Println("[P]rint the dataset as JSON")
	fmt.Println("[R]eset the dataset")
	fmt.Println("[Q]uit")
	choice, _ := getInput("Enter your choice: ", reader)

	switch strings.ToLower(choice) {
	case "a":
		data.WriteFiles()
	case "c":
		data.CheckIfCaptionsExist()
	case "m":
		data.MoveCaptionsToImages()
	case "p":
		byteArr, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(byteArr))
	case "r":
		data.Images = make(map[string]Image)
	case "q":
		return
	default:
		fmt.Println("Invalid choice")
	}

	data.promptOption()
}

func getInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Println(prompt)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}
