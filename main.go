package main

import (
	"find-image-tag/tui"
	"fmt"
	"github.com/nokusukun/roggy"
	"github.com/nokusukun/stemp"
	"strings"
)

var roggyPrinter = roggy.Printer("main-service")
var roggyNoTrace = roggy.Printer("main-service")
var overwrite bool

func main() {
	//// Read all the filenames of image files in the image directory
	//
	//roggy.LogLevel = roggy.TypeDebug
	//roggyPrinter.Debug("Starting main service")
	//roggyPrinter.Sync = true
	//roggyNoTrace.NoTrace = true
	//
	//data := InitDataSet()
	//
	//data.promptOption()

	roggy.LogLevel = roggy.TypeRoggy
	tui.Main()
}

var validChoice = true

//func (data *DataSet) promptOption() {
//
//	reader := bufio.NewReader(os.Stdin)
//
//	roggy.Flush()
//
//	countImagesWithCaptions := data.countImagesWithCaptions()
//	countCaptionDirectoryMatchImageDirectory := data.countCaptionDirectoryMatchImageDirectory()
//	countImagesWithoutCaptions := data.countImagesWithoutCaptions()
//	countPending := data.countPending()
//	countFiles := data.countFiles()
//	countImages := data.countImages()
//	countOverwrites := data.countOverwrites(overwrite)
//	countCaptionsToMerge := data.countCaptionsToMerge()
//	countTotalCaptions := data.countPending() + data.countImagesWithCaptions()
//	countImagesWithCaptionsNextToThem := data.countImagesWithCaptionsNextToThem()
//	offSet := len(fmt.Sprintf(" + %d -> %d", countImagesWithCaptionsNextToThem, countImagesWithCaptions))
//	moveString := fmt.Sprintf("{0:w=%d,j=r} + {1:j=r} -> {2:j=r}", 30-offSet)
//	//countToMove := data.countImagesWithCaptions() - data.countCaptionDirectoryMatchImageDirectory()
//	nul := ""
//	var ow any
//
//	var overwriteString string
//	var moveOverwriteString string
//	if overwrite {
//		overwriteString = "[O]verwriting existing caption files even if they exist"
//		moveOverwriteString = "and overwrite captions to the image files"
//		ow = countOverwrites
//	} else {
//		overwriteString = "[O]nly moving/hardlinking caption files that don't exist"
//		moveOverwriteString = "missing captions to the image files"
//		ow = nul
//	}
//
//	toPrint := []string{
//		"1::" + roggy.Rainbowize("---") + " Stats " + roggy.Rainbowize("---"),
//		"2::",
//		stemp.Inline("2::{0:w=30,j=r} | Images with captions", countImagesWithCaptions),
//		stemp.Inline("2::{0:w=30,j=r} | Images with captions that match directories", countCaptionDirectoryMatchImageDirectory),
//		stemp.Inline("2::{0:w=30,j=r} | Missing captions", countImagesWithoutCaptions),
//		stemp.Inline("2::{0:w=30,j=r} | Pending text files", countPending),
//		"1::" + roggy.Rainbowize("---") + " Image Captioning " + roggy.Rainbowize("---"),
//		"2::",
//		stemp.Inline("2::{0:w=30,j=r} | [+] Add files to the dataset", countFiles),
//		stemp.Inline("2::{0:w=30,j=r} | [+c] Add captions to the dataset", countTotalCaptions),
//		stemp.Inline("2::{0:w=30,j=r} | [+i] Add images to the dataset", countImages),
//		stemp.Inline("2::{0:w=30,j=r} | [C]heck if each image has a caption", countImages),
//		stemp.Inline("2::{0:w=30,j=r} | [P]rint the dataset as JSON", nul),
//		stemp.Inline("2::{0:w=30,j=r} | [R]eset the dataset", nul),
//		stemp.Inline("2::{0:w=30,j=r} | [W]rite the dataset as a JSON file", nul),
//		stemp.Inline("2::{0:w=30,j=r} | Append [t]ext files to matching images", countPending),
//		stemp.Inline("2::{0:w=30,j=r} | Check for captions without matching [i]mages", nul),
//		stemp.Inline("2::{0:w=30,j=r} | [Q]uit", nul),
//		"1::" + roggy.Rainbowize("---") + " Actions " + roggy.Rainbowize("---"),
//		"2::",
//		stemp.Inline("2::{0:w=30,j=r} | {1} | {2:w=10,j=r}", ow, overwrite, overwriteString),
//		stemp.Inline("2::"+moveString+" | [Move] {3}", countImagesWithCaptionsNextToThem, countOverwrites, countImagesWithCaptions, moveOverwriteString),
//		stemp.Inline("2::"+moveString+" | [Hardlink] {3}", countImagesWithCaptionsNextToThem, countOverwrites, countImagesWithCaptions, moveOverwriteString),
//		stemp.Inline("2::{0:w=30,j=r} | [Merge] new captions to existing captions", countCaptionsToMerge),
//		stemp.Inline("2::{0:w=30,j=r} | [Append] new tags to captions (dir)", countImagesWithCaptionsNextToThem),
//		stemp.Inline("2::{0:w=30,j=r} | Replace spaces with [_]", nul),
//	}
//
//	printLogs(toPrint)
//
//	if !validChoice {
//		roggyPrinter.Errorf("Invalid choice. Please try again.")
//		validChoice = true
//	}
//	choice, _ := getInput("Enter your choice: ", reader)
//
//	choice = strings.ToLower(choice)
//
//	var directory string
//	switch choice {
//	case "+", "+c", "+i":
//		fmt.Print("Enter the directory to read: ")
//		_, _ = fmt.Scanln(&directory)
//		if _, err := os.Stat(directory); os.IsNotExist(err) {
//			choice = ""
//		}
//	}
//
//	switch choice {
//	case "+":
//		data.WriteFiles(addBoth, directory)
//	case "+c":
//		data.WriteFiles(addCaption, directory)
//	case "+i":
//		data.WriteFiles(addImage, directory)
//	case "c":
//		data.CheckIfCaptionsExist()
//	case "o":
//		overwrite = !overwrite
//	case "move":
//		data.CaptionsToImages(move, overwrite)
//	case "hardlink":
//		data.CaptionsToImages(hardlink, overwrite)
//	case "p":
//		data.prettyJson()
//	case "w":
//		data.writeJson()
//	case "r":
//		*data = *InitDataSet()
//	case "t":
//		data.AppendCaptionsConcurrently()
//	case "i":
//		data.checkForMissingImages()
//	case "merge":
//		data.CaptionsToImages(merge, overwrite)
//	case "append":
//		appendNewTags()
//	case "_":
//		data.ReplaceSpaces()
//	case "q":
//		return
//	default:
//		validChoice = false
//	}
//
//	data.promptOption()
//}
//
//func getInput(prompt string, reader *bufio.Reader) (string, error) {
//	roggyPrinter.Infof(prompt)
//	input, err := reader.ReadString('\n')
//	return strings.TrimSpace(input), err
//}

func printLogs(toPrint []string) {
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
			printLog(bufferLevel, output)
			buffer = []string{}
		}

		bufferLevel = debugLevel
		buffer = append(buffer, dbgString[1])
	}

	// last batch
	if len(buffer) > 0 {
		output := strings.Join(buffer, "\n")
		printLog(bufferLevel, output)
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

func printLog(debugLevel string, output string) {
	nilMap := map[string]any{}
	if funcToUse, ok := logFunctions[debugLevel]; ok {
		toOutput := stemp.Compile(output, nilMap)
		funcToUse(toOutput)
	} else {
		fmt.Printf("Debug level outside valid range: %s", debugLevel)
		return
	}
}
