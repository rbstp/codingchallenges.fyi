package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	lineCount, wordCount, byteCount, charCount := 0, 0, 0, 0
	var err error

	if countLinesFlag {
		if useStdin {
			lineCount, err = countLinesFromStdin()
		} else {
			lineCount, err = countLines(filename)
		}
		handleError(err)
	}

	if countWordsFlag {
		if useStdin {
			wordCount, err = countWordsFromStdin()
		} else {
			wordCount, err = countWords(filename)
		}
		handleError(err)
	}

	if countBytesFlag {
		if useStdin {
			byteCount, err = countBytesFromStdin()
		} else {
			byteCount, err = countBytes(filename)
		}
		handleError(err)
	}

	if countCharsFlag {
		if useStdin {
			charCount, err = countCharactersFromStdin()
		} else {
			charCount, err = countCharacters(filename)
		}
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

func countBytesFromStdin() (int, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return 0, fmt.Errorf("error reading from stdin: %v", err)
	}
	return len(content), nil
}

func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	return countLinesFromReader(file)
}

func countLinesFromStdin() (int, error) {
	return countLinesFromReader(os.Stdin)
}

func countLinesFromReader(reader io.Reader) (int, error) {
	lineCount := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading: %v", err)
	}
	return lineCount, nil
}

func countWords(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	return countWordsFromReader(file)
}

func countWordsFromStdin() (int, error) {
	return countWordsFromReader(os.Stdin)
}

func countWordsFromReader(reader io.Reader) (int, error) {
	wordCount := 0
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading: %v", err)
	}
	return wordCount, nil
}

func countCharacters(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	return countCharactersFromReader(file)
}

func countCharactersFromStdin() (int, error) {
	return countCharactersFromReader(os.Stdin)
}

func countCharactersFromReader(reader io.Reader) (int, error) {
	charCount := 0
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		charCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading: %v", err)
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
	cat test.txt | ccwc -l`)
}
