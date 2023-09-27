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
var overwrite bool

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
		"2::",
		"2::{countImagesWithCaptions:w=30,j=r} | Images with captions",
		"2::{countCaptionDirectoryMatchImageDirectory:w=30,j=r} | Images with captions that match directories",
		"2::{countImagesWithoutCaptions:w=30,j=r} | Missing captions",
		"2::{countPending:w=30,j=r} | Pending text files",
		"1::" + roggy.Rainbowize("---") + " Image Captioning " + roggy.Rainbowize("---"),
		"2::",
		"2::{countFiles:w=30,j=r} | [+] Add files to the dataset",
		"2::{countImages:w=30,j=r} | [C]heck each image has a caption",
		"2::{nul:w=30,j=r} | [P]rint the dataset as JSON",
		"2::{nul:w=30,j=r} | [R]eset the dataset",
		"2::{nul:w=30,j=r} | [W]rite the dataset as a JSON file",
		"2::{countPending:w=30,j=r} | Append [t]ext files to matching images",
		"2::{nul:w=30,j=r} | Check for captions without matching [i]mages",
		"2::{nul:w=30,j=r} | [Q]uit",
		"1::" + roggy.Rainbowize("---") + " Actions " + roggy.Rainbowize("---"),
		"2::",
		"2::{nul:w=30,j=r} | {overwrite} | {overwriteString:w=10,j=r}",
		"2::{countImagesWithoutCaptions:w=30,j=r} | [Move] captions to the image files",     // TODO: Only count overwrites if overwrite is true
		"2::{countImagesWithoutCaptions:w=30,j=r} | [Hardlink] captions to the image files", // TODO: Use countOverwrites() after fixing implementation
		"2::{countExistingCaptions:w=30,j=r} | [Merge] captions to the image files",         // TODO: Fix countExistingCaptions() implementation
		"2::{countExistingCaptions:w=30,j=r} | [Append] new tags to the caption file",
		"2::{nul:w=30,j=r} | Replace spaces with [_]",
	}

	var overwriteString string

	if overwrite {
		overwriteString = "[O]verwriting existing caption files even if they exist"
	} else {
		overwriteString = "[O]nly moving/hardlinking caption files that don't exist"
	}

	values := map[string]any{
		"countImagesWithCaptions":                  data.countImagesWithCaptions(),
		"countCaptionDirectoryMatchImageDirectory": data.countCaptionDirectoryMatchImageDirectory(),
		"countImagesWithoutCaptions":               data.countImagesWithoutCaptions(),
		"countPending":                             data.countPending(),
		"countFiles":                               data.countFiles(),
		"countImages":                              data.countImages(),
		"countOverwrites":                          data.countOverwrites(),       // pass overwrite bool
		"countExistingCaptions":                    data.countExistingCaptions(), // pass overwrite bool
		"nul":                                      "",
		"overwrite":                                overwrite,
		"overwriteString":                          overwriteString,
	}

	printLogs(toPrint, values)

	choice, _ := getInput("Enter your choice: ", reader)

	switch strings.ToLower(choice) {
	case "+":
		data.WriteFiles()
	case "c":
		data.CheckIfCaptionsExist()
	case "o":
		overwrite = !overwrite
	case "move":
		data.CaptionsToImages(move)
	case "hardlink":
		data.CaptionsToImages(hardlink)
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
	case "merge":
		data.CaptionsToImages(merge)
	case "append":
		data.appendNewTags()
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
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	bytes, _ := json.MarshalIndent(data.Images, "", "  ")
	_, _ = file.Write(bytes)
}

func printLogs(toPrint []string, values map[string]any) {
	var buffer []string
	bufferLevel := ""

	for _, p := range toPrint {
		if p == "" {
			roggyNoTrace.Infof(strings.Repeat(" ", 60))
			continue
		}

		dbgString := strings.Split(p, "::")
		debugLevel := dbgString[0]

		if dbgString[1] == "" {
			dbgString[1] = strings.Repeat(" ", 60)
		}

		if bufferLevel != "" && bufferLevel != debugLevel {
			output := strings.Join(buffer, "\n")
			printLog(bufferLevel, output, values)
			buffer = []string{}
		}

		bufferLevel = debugLevel
		buffer = append(buffer, dbgString[1])
	}

	// last batch
	if len(buffer) > 0 {
		output := strings.Join(buffer, "\n")
		printLog(bufferLevel, output, values)
	}
}

var logFunctions = map[string]func(f string, message ...interface{}){
	//"-1": roggyPrinter.Roggyf,
	"0": roggyPrinter.Errorf,
	"1": roggyPrinter.Noticef,
	"2": roggyPrinter.Infof,
	"3": roggyPrinter.Verbosef,
	"4": roggyPrinter.Debugf,
}

func printLog(debugLevel string, output string, values map[string]any) {
	if logFunc, ok := logFunctions[debugLevel]; ok {
		toOutput := stemp.Compile(output, values)
		logFunc(toOutput)
	} else {
		fmt.Printf("Debug level outside valid range: %s", debugLevel)
		return
	}
}
