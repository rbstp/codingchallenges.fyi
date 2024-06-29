package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(`Usage: ccwc [option] <file>
		Options:
			-c  Count the number of bytes in the file.
			-l  Count the number of lines in the file.
			-w  Count the number of words in the file.
			-m  Count the number of characters in the file.
		Example:
			ccwc -c myfile.txt`)
		os.Exit(1)
	}

	option := os.Args[1]
	filename := os.Args[2]

	switch option {
	case "-c":
		byteCount, err := countBytes(filename)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("%d %s\n", byteCount, filename)
	case "-l":
		lineCount, err := countLines(filename)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("%d %s\n", lineCount, filename)
	case "-w":
		wordCount, err := countWords(filename)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("%d %s\n", wordCount, filename)
	case "-m":
		charCount, err := countCharacters(filename)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("%d %s\n", charCount, filename)
	default:
		fmt.Println(`Invalid option. Available options are:
			-c  Count the number of bytes in the file.
			-l  Count the number of lines in the file.
			-w  Count the number of words in the file.
			-m  Count the number of characters in the file.
		Please check the usage instructions for more information.`)
		os.Exit(1)
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

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
