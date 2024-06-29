package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[len(os.Args)-1] // Last argument is the filename

	if len(os.Args) == 2 { // No option provided, default to -c, -l, -w
		handleCounts(filename, true, true, true, false)
		return
	}

	option := os.Args[1]
	switch option {
	case "-c":
		handleCounts(filename, false, false, true, false)
	case "-l":
		handleCounts(filename, true, false, false, false)
	case "-w":
		handleCounts(filename, false, true, false, false)
	case "-m":
		handleCounts(filename, false, false, false, true)
	default:
		fmt.Println("Invalid option.")
		printUsage()
		os.Exit(1)
	}
}

func handleCounts(filename string, countLinesFlag, countWordsFlag, countBytesFlag, countCharsFlag bool) {
	lineCount, wordCount, byteCount, charCount := 0, 0, 0, 0
	var err error

	if countLinesFlag {
		lineCount, err = countLines(filename)
		handleError(err)
	}

	if countWordsFlag {
		wordCount, err = countWords(filename)
		handleError(err)
	}

	if countBytesFlag {
		byteCount, err = countBytes(filename)
		handleError(err)
	}

	if countCharsFlag {
		charCount, err = countCharacters(filename)
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

func countBytes(filename string) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filename, err)
	}
	return len(content), nil
}

func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filename, err)
	}
	return lineCount, nil
}

func countWords(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filename, err)
	}
	return wordCount, nil
}

func countCharacters(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	charCount := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes) // Split by runes (characters)
	for scanner.Scan() {
		charCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filename, err)
	}
	return charCount, nil
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
	ccwc -c myfile.txt`)
}
