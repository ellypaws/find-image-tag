package tui

import (
	"find-image-tag/tui/sender"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const desiredSteps = 200

func addMultiple(steps int) tea.Cmd {
	if steps == 0 {
		steps = desiredSteps
	}
	return func() tea.Msg {
		return addMultipleMsg{current: 1, total: steps} // Start with the first step and desired total steps
	}
}

func (m model) singleDirToMultiple(filter int, directory string) tea.Cmd {
	dirSliceFullPath, err := checkSubdirectoriesRecursively(directory, 10)
	if err != nil {
		// Handle the error accordingly. For now, returning nil.
		return nil
	}

	msg := func() tea.Msg {
		return writeFilesMsg{
			current: 1,
			total:   len(dirSliceFullPath),
			filter:  filter,
			dirs:    dirSliceFullPath,
			wg:      &sync.WaitGroup{},
		}
	}

	var addToLog []tea.Cmd
	var mutex sync.Mutex
	wg := &sync.WaitGroup{}

	for _, dir := range dirSliceFullPath {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			currentTime := time.Now()
			mutex.Lock()
			addToLog = append(addToLog, func() tea.Msg {
				return sender.ResultMsg{
					Food:     dir,
					Duration: time.Since(currentTime),
				}
			})
			mutex.Unlock()
		}(dir)
	}

	wg.Wait()

	return tea.Batch(msg, tea.Batch(addToLog...)) // Adjusted for Bubbletea's Batch usage.
}

func checkSubdirectoriesRecursively(directory string, maxChildren int) ([]string, error) {
	var result []string

	// Recursive helper function
	var processDir func(currentDir string) error
	processDir = func(currentDir string) error {
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return err
		}

		var subDirs []string
		for _, entry := range entries {
			if entry.IsDir() {
				subDirPath := filepath.Join(currentDir, entry.Name())
				subDirs = append(subDirs, subDirPath)
			}
		}

		if len(subDirs) > maxChildren {
			// If more than maxChildren, recursively process each subdirectory
			for _, subDir := range subDirs {
				if err := processDir(subDir); err != nil {
					return err
				}
			}
		} else {
			// Otherwise, simply add the current directory to the result
			result = append(result, currentDir)
		}

		return nil
	}

	// Begin processing from the given directory
	if err := processDir(directory); err != nil {
		return nil, err
	}

	return result, nil
}

func processDirectory(m *model, filter int, currentMsg writeFilesMsg) tea.Cmd {
	currentMsg.wg.Add(1)
	go m.DataSet.WriteFiles(filter, currentMsg.dirs[0])
	sendAgain := func(t time.Time) tea.Msg {
		return writeFilesMsg{
			current: currentMsg.current + 1,
			total:   currentMsg.total,
			filter:  currentMsg.filter,
			dirs:    currentMsg.dirs[1:],
			wg:      currentMsg.wg,
		}
	}
	return tea.Tick(time.Millisecond, sendAgain)
}

func formatWithComma(i int) string {
	in := strconv.Itoa(i)
	var out strings.Builder
	l := len(in)
	for i, r := range in {
		_, _ = out.WriteRune(r)
		if (l-i-1)%3 == 0 && i < l-1 {
			_, _ = out.WriteRune(',')
		}
	}
	return out.String()
}
