package surf

import (
	"bufio"
	"os"

	"github.com/schollz/progressbar/v3"
)

func readLines(path string) ([]string, error) {
	// read the file into a string slice
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// iterate over each line in the file
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func lineCounter(filename string) (int64, error) {
	// count the lines in a file
	// return the line count
	// and an error if any
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// create a buffer to store the file contents
	buf := make([]byte, 1024)
	var count int64
	count = 0
	lineSep := []byte{'\n'}

	for {
		// read a chunk of the file into the buffer
		c, err := file.Read(buf)
		if err != nil {
			break
		}

		// iterate over each line in the buffer
		for _, b := range buf[:c] {
			if b == lineSep[0] {
				count++
			}
		}
	}

	return count, nil
}

// initialise and display progress bar
func initBar(totalLines int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		totalLines,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription("[cyan]Running..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	return bar
}

// add 1 to the progress bar
func incrementBar(bar *progressbar.ProgressBar, incrementCount int) {
	bar.Add(incrementCount)
}
