package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	globalIgnoreFile = filepath.Join(os.Getenv("HOME"), ".tarrareignore")
	globalIgnoreList []string
	outputFileName   string
)

func main() {
	directory := flag.String("dir", "", "Directory to process (default: current directory)")
	outputFile := flag.String("output", "", "Output file name (default: current directory name)")
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of files/directories to ignore")
	flag.Parse()

	// Handle default directory
	if *directory == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}
		*directory = currentDir
	}

	loadGlobalIgnoreList()
	ignoreList := append(globalIgnoreList, strings.Split(*ignoreFlag, ",")...)
	if *outputFile == "" {
		outputFileName = filepath.Base(*directory) + ".txt"
	} else {
		outputFileName = *outputFile
	}

	// Add the output file to the ignore list
	ignoreList = append(ignoreList, outputFileName)

	out, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	fmt.Fprintf(out, "Root Directory: %s\n\n", filepath.Base(*directory))

	err = processDirectory(*directory, *directory, ignoreList, out)
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		return
	}

	fmt.Printf("Successfully created %s\n", *outputFile)
}
func loadGlobalIgnoreList() {
	file, err := os.Open(globalIgnoreFile)
	if err != nil {
		return // It's okay if the file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		globalIgnoreList = append(globalIgnoreList, strings.TrimSpace(scanner.Text()))
	}
}

func processDirectory(baseDir, currentDir string, ignoreList []string, out *os.File) error {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", currentDir, err)
	}

	for _, entry := range entries {
		relativePath, err := filepath.Rel(baseDir, filepath.Join(currentDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %v", entry.Name(), err)
		}

		if shouldIgnore(relativePath, ignoreList) {
			continue
		}

		if entry.IsDir() {
			err = processDirectory(baseDir, filepath.Join(currentDir, entry.Name()), ignoreList, out)
			if err != nil {
				fmt.Printf("Warning: Error processing directory %s: %v\n", relativePath, err)
				continue
			}
		} else {
			err = processFile(baseDir, filepath.Join(currentDir, entry.Name()), out)
			if err != nil {
				fmt.Printf("Warning: Error processing file %s: %v\n", relativePath, err)
				continue
			}
		}
	}

	return nil
}

func processFile(baseDir, filePath string, out *os.File) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("not a regular file")
	}

	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return fmt.Errorf("error getting relative path: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	fmt.Fprintf(out, "--- BEGIN FILE: %s ---\n", relativePath)
	fmt.Fprintf(out, "%s\n", string(content))
	fmt.Fprintf(out, "--- END FILE: %s ---\n\n", relativePath)

	return nil
}

func shouldIgnore(path string, ignoreList []string) bool {
	for _, ignore := range ignoreList {
		if ignore == "" {
			continue
		}
		if matched, _ := filepath.Match(ignore, filepath.Base(path)); matched {
			return true
		}
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}
