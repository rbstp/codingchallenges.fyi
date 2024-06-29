package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ccwc -c <file> or ccwc -l <file>")
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
	default:
		fmt.Println("Invalid option. Use -c for byte count or -l for line count.")
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

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filename, err)
	}
	return lineCount, nil
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
