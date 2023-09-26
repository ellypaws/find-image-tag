package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/TylerBrock/colorjson"
	"github.com/nokusukun/roggy"
	"github.com/nokusukun/stemp"
	"os"
	"strings"
)

var roggyPrinter = roggy.Printer("main-service")
var roggyNoTrace = roggy.Printer("main-service")

func main() {
	// Read all the filenames of image files in the image directory

	roggy.LogLevel = roggy.TypeDebug
	roggyPrinter.Debug("Starting main service")
	roggyPrinter.Sync = true
	roggyNoTrace.NoTrace = true

	data := InitDataSet()

	data.promptOption()
}

func (data *DataSet) promptOption() {
	reader := bufio.NewReader(os.Stdin)

	roggy.Flush()

	toPrint := []string{
		"1::" + roggy.Rainbowize("---") + " Stats " + roggy.Rainbowize("---"),
		"2::{countImagesWithCaptions:w=5,j=r} | Images with captions",
		"2::{countCaptionDirectoryMatchImageDirectory:w=5,j=r} | Images with captions that match directories",
		"2::{countImagesWithoutCaptions:w=5,j=r} | Missing captions",
		"2::{countPending:w=5,j=r} | Pending text files",
		"",
		"1::" + roggy.Rainbowize("---") + " Image Captioning " + roggy.Rainbowize("---"),
		"2::{countFiles:w=5,j=r} | [A]dd files to the dataset",
		"2::{countImages:w=5,j=r} | [C]heck each image has a caption",
		"2::{nul:w=5,j=r} | [P]rint the dataset as JSON",
		"2::{nul:w=5,j=r} | [R]eset the dataset",
		"2::{nul:w=5,j=r} | [W]rite the dataset as a JSON file",
		"2::{countPending:w=5,j=r} | Append [t]ext files to matching images",
		"2::{nul:w=5,j=r} | Check for captions without matching [i]mages",
		"2::{nul:w=5,j=r} | [Q]uit",
		"",
		"1::" + roggy.Rainbowize("---") + " Actions " + roggy.Rainbowize("---"),
		"2::{countImagesWithCaptions:w=5,j=r} | [M]ove captions to the image files",
		"2::{countImagesWithCaptions:w=5,j=r} | C[o]py captions to the image files",
		"2::{nul:w=5,j=r} | Replace spaces with [_]",
	}

	values := map[string]any{
		"countImagesWithCaptions":                  data.countImagesWithCaptions(),
		"countCaptionDirectoryMatchImageDirectory": data.countCaptionDirectoryMatchImageDirectory(),
		"countImagesWithoutCaptions":               data.countImagesWithoutCaptions(),
		"countPending":                             data.countPending(),
		"countFiles":                               data.countFiles(),
		"countImages":                              data.countImages(),
		"nul":                                      "",
	}

	for _, p := range toPrint {
		if p == "" {
			roggyNoTrace.Infof(strings.Repeat(" ", 60))
			continue
		}
		dbgString := strings.Split(p, "::")
		debugLevel := dbgString[0]
		if debugLevel == "1" {
			roggyPrinter.Noticef(dbgString[1])
			continue
		}
		if debugLevel == "2" {
			p = dbgString[1]
		}
		roggyPrinter.Infof(stemp.Compile(p, values))
	}

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
		data.appendCaptionsConcurrently()
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

var logFunctions = map[string]func(f string, message ...interface{}){
	//"-1": roggyPrinter.Roggyf,
	"0": roggyPrinter.Errorf,
	"1": roggyPrinter.Noticef,
	"2": roggyPrinter.Infof,
	"3": roggyPrinter.Verbosef,
	"4": roggyPrinter.Debugf,
}

func printLog(debugLevel string, output string) {
	if logFunc, ok := logFunctions[debugLevel]; ok {
		logFunc(output)
	} else {
		fmt.Printf("Debug level outside valid range: %s", debugLevel)
		return
	}
}
