package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCountBytes(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer file.Close()

	byteCount, err := countBytes(file)
	if err != nil {
		t.Errorf("countBytes failed: %v", err)
	}

	expectedByteCount := 342190
	if byteCount != expectedByteCount {
		t.Errorf("expected %d bytes, got %d", expectedByteCount, byteCount)
	}
}

func TestCountLines(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer file.Close()

	lineCount, _, _, err := countLinesWordsChars(file, true, false, false)
	if err != nil {
		t.Errorf("countLines failed: %v", err)
	}

	expectedLineCount := 7145
	if lineCount != expectedLineCount {
		t.Errorf("expected %d lines, got %d", expectedLineCount, lineCount)
	}
}

func TestCountWords(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer file.Close()

	_, wordCount, _, err := countLinesWordsChars(file, false, true, false)
	if err != nil {
		t.Errorf("countWords failed: %v", err)
	}

	expectedWordCount := 58164
	if wordCount != expectedWordCount {
		t.Errorf("expected %d words, got %d", expectedWordCount, wordCount)
	}
}

func TestCountCharacters(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer file.Close()

	_, _, charCount, err := countLinesWordsChars(file, false, false, true)
	if err != nil {
		t.Errorf("countCharacters failed: %v", err)
	}

	expectedCharCount := 325002
	if charCount != expectedCharCount {
		t.Errorf("expected %d characters, got %d", expectedCharCount, charCount)
	}
}

func TestCountFromStdin(t *testing.T) {
	input := "Hello world\nThis is a test\n"
	expectedLineCount := 2
	expectedWordCount := 6
	expectedCharCount := 25

	reader := strings.NewReader(input)

	lineCount, wordCount, charCount, err := countLinesWordsChars(reader, true, true, true)
	if err != nil {
		t.Errorf("count from stdin failed: %v", err)
	}

	if lineCount != expectedLineCount {
		t.Errorf("expected %d lines, got %d", expectedLineCount, lineCount)
	}
	if wordCount != expectedWordCount {
		t.Errorf("expected %d words, got %d", expectedWordCount, wordCount)
	}
	if charCount != expectedCharCount {
		t.Errorf("expected %d characters, got %d", expectedCharCount, charCount)
	}
}

func TestHandleCounts(t *testing.T) {
	var buf bytes.Buffer
	input := "Hello world\nThis is a test\n"
	expectedLineCount := 2
	expectedWordCount := 6
	expectedByteCount := len(input)

	reader := strings.NewReader(input)

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCounts(reader, "", true, true, true, false, true)

	w.Close()
	os.Stdout = stdout

	io.Copy(&buf, r)

	output := buf.String()
	expectedOutput := fmt.Sprintf("%8d %8d %8d %s\n", expectedLineCount, expectedWordCount, expectedByteCount, "")

	if output != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, output)
	}
}
