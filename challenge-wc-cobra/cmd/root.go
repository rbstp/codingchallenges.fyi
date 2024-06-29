package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	countLinesFlag bool
	countWordsFlag bool
	countBytesFlag bool
	countCharsFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "challenge-wc-cobra [file]",
	Short: "challenge-wc-cobra is a CLI tool to count lines, words, characters, and bytes in a file",
	Long: `A fast and flexible command line tool built in Go to count lines, words,
characters, and bytes in a file or standard input.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := ""
		if len(args) > 0 {
			filename = args[0]
		}
		if !countLinesFlag && !countWordsFlag && !countBytesFlag && !countCharsFlag {
			countLinesFlag = true
			countWordsFlag = true
			countBytesFlag = true
		}
		handleCounts(os.Stdin, filename, countLinesFlag, countWordsFlag, countBytesFlag, countCharsFlag, filename == "")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&countLinesFlag, "lines", "l", false, "Count lines")
	rootCmd.Flags().BoolVarP(&countWordsFlag, "words", "w", false, "Count words")
	rootCmd.Flags().BoolVarP(&countBytesFlag, "bytes", "c", false, "Count bytes")
	rootCmd.Flags().BoolVarP(&countCharsFlag, "chars", "m", false, "Count characters")
}

func handleCounts(stdin io.Reader, filename string, countLinesFlag, countWordsFlag, countBytesFlag, countCharsFlag bool, useStdin bool) {
	var reader io.Reader
	if useStdin {
		reader = stdin
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
