package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		os.Exit(1)
	}

	var filename string
	var useStdin bool

	if len(os.Args) == 2 {
		if os.Args[1][0] == '-' {
			useStdin = true
		} else {
			filename = os.Args[1]
		}
	} else {
		filename = os.Args[2]
	}

	if len(os.Args) == 2 && filename != "" { // No option provided, default to -c, -l, -w
		handleCounts(filename, true, true, true, false, useStdin)
		return
	}

	option := os.Args[1]
	switch option {
	case "-c":
		handleCounts(filename, false, false, true, false, useStdin)
	case "-l":
		handleCounts(filename, true, false, false, false, useStdin)
	case "-w":
		handleCounts(filename, false, true, false, false, useStdin)
	case "-m":
		handleCounts(filename, false, false, false, true, useStdin)
	default:
		fmt.Println("Invalid option.")
		printUsage()
		os.Exit(1)
	}
}

func handleCounts(filename string, countLinesFlag, countWordsFlag, countBytesFlag, countCharsFlag bool, useStdin bool) {
	var reader io.Reader
	if useStdin {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			handleError(fmt.Errorf("error opening file %s: %v", filename, err))
		}
		defer file.Close()
		reader = file
	}

	lineCount, wordCount, byteCount, charCount := 0, 0, 0, 0
	var err error

	if countLinesFlag || countWordsFlag || countCharsFlag {
		lineCount, wordCount, charCount, err = countLinesWordsChars(reader, countLinesFlag, countWordsFlag, countCharsFlag)
		handleError(err)
	}

	if countBytesFlag {
		if seeker, ok := reader.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}
		byteCount, err = countBytes(reader)
		handleError(err)
	}

	if countLinesFlag && countWordsFlag && countBytesFlag {
		fmt.Printf("%8d %8d %8d %s\n", lineCount, wordCount, byteCount, filename)
	} else if countLinesFlag {
		fmt.Printf("%8d %s\n", lineCount, filename)
	} else if countWordsFlag {
		fmt.Printf("%8d %s\n", wordCount, filename)
	} else if countBytesFlag {
		fmt.Printf("%8d %s\n", byteCount, filename)
	} else if countCharsFlag {
		fmt.Printf("%8d %s\n", charCount, filename)
	}
}

func countBytes(reader io.Reader) (int, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return 0, fmt.Errorf("error reading: %v", err)
	}
	return len(content), nil
}

func countLinesWordsChars(reader io.Reader, countLinesFlag, countWordsFlag, countCharsFlag bool) (int, int, int, error) {
	lineCount, wordCount, charCount := 0, 0, 0
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if countLinesFlag {
			lineCount++
		}
		if countWordsFlag || countCharsFlag {
			text := scanner.Text()
			if countWordsFlag {
				wordCount += len(strings.Fields(text))
			}
			if countCharsFlag {
				charCount += len([]rune(text))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, 0, fmt.Errorf("error reading: %v", err)
	}

	return lineCount, wordCount, charCount, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: ccwc [option] <file>
Options:
	-c  Count the number of bytes in the file.
	-l  Count the number of lines in the file.
	-w  Count the number of words in the file.
	-m  Count the number of characters in the file.
Example:
	cat test.txt | ccwc -l`)
}
