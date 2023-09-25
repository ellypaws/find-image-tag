package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/TylerBrock/colorjson"
	"github.com/nokusukun/roggy"
	"os"
	"strings"
)

var roggyPrinter = roggy.Printer("main-service")

func main() {
	// Read all the filenames of image files in the image directory

	roggy.LogLevel = roggy.TypeDebug
	roggyPrinter.Debug("Starting main service")

	data := InitDataSet()

	data.promptOption()
}

func (data *DataSet) promptOption() {
	reader := bufio.NewReader(os.Stdin)

	roggy.Flush()
	roggyPrinter.Infof("--- Stats ---")
	roggyPrinter.Infof("Images with captions: %i", data.countImagesWithCaptions())
	roggyPrinter.Infof("Missing captions: %i", data.countImagesWithoutCaptions())
	roggyPrinter.Infof("Pending text files: %i", data.countPending())
	roggyPrinter.Infof("--- Image Captioning ---")
	roggyPrinter.Infof("[%i][A]dd files to the dataset", data.countFiles())
	roggyPrinter.Infof("[%i][C]heck if captions exist", data.countImages())
	roggyPrinter.Infof("[P]rint the dataset as JSON")
	roggyPrinter.Infof("[R]eset the dataset")
	roggyPrinter.Infof("[W]rite the dataset as a JSON file")
	roggyPrinter.Infof("[%i]Append [t]ext files to matching images", data.countPending())
	roggyPrinter.Infof("Check for captions without matching [i]mages")
	roggyPrinter.Infof("[Q]uit")
	roggyPrinter.Infof("")
	roggyPrinter.Infof("[%i][M]ove captions to the image files", data.countImagesWithCaptions())
	roggyPrinter.Infof("[%i]C[o]py captions to the image files", data.countImagesWithCaptions())
	roggyPrinter.Infof("Replace spaces with [_]")
	choice, _ := getInput("Enter your choice: ", reader)

	switch strings.ToLower(choice) {
	case "a":
		data.WriteFiles()
	case "c":
		data.CheckIfCaptionsExist()
	case "m":
		data.CaptionsToImages(true)
	case "o":
		data.CaptionsToImages(false)
	case "p":
		data.prettyJson()
	case "w":
		data.writeJson()
	case "r":
		*data = *InitDataSet()
	case "t":
		data.appendCaptions()
	case "i":
		data.checkForMissingImages()
	case "_":
		data.replaceSpaces()
	case "q":
		return
	default:
		fmt.Println("Invalid choice")
	}

	data.promptOption()
}

func getInput(prompt string, reader *bufio.Reader) (string, error) {
	roggyPrinter.Infof(prompt)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func (data *DataSet) prettyJson() {
	var obj map[string]any
	bytes, _ := json.Marshal(data)
	_ = json.Unmarshal(bytes, &obj)

	formatter := colorjson.NewFormatter()
	formatter.Indent = 2

	byteArray, err := formatter.Marshal(obj)
	roggyPrinter.Debugf("Byte array: %s", byteArray)
	roggyPrinter.Debugf("Error: %s", err)
	roggyPrinter.Infof(string(byteArray))
}

func (data *DataSet) writeJson() {
	roggyPrinter.Infof("Writing dataset to file...")
	file, _ := os.Create("dataset.json")
	defer file.Close()

	bytes, _ := json.MarshalIndent(data.Images, "", "  ")
	_, _ = file.Write(bytes)
}
