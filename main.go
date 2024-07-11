package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	globalIgnoreFile = filepath.Join(os.Getenv("HOME"), ".tarrareignore")
	globalIgnoreList []string
)

func main() {
	directory := flag.String("dir", ".", "Directory to process")
	outputFile := flag.String("output", "", "Output file name (default: current directory name)")
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of files/directories to ignore")
	flag.Parse()

	loadGlobalIgnoreList()
	ignoreList := append(globalIgnoreList, strings.Split(*ignoreFlag, ",")...)

	if *outputFile == "" {
		*outputFile = filepath.Base(*directory) + ".txt"
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	fmt.Fprintf(out, "Root Directory: %s\n\n", filepath.Base(*directory))

	err = processDirectory(*directory, *directory, ignoreList, out, 0)
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

func processDirectory(baseDir, currentDir string, ignoreList []string, out *os.File, depth int) error {
	entries, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		relativePath, err := filepath.Rel(baseDir, filepath.Join(currentDir, entry.Name()))
		if err != nil {
			return err
		}

		if shouldIgnore(relativePath, ignoreList) {
			continue
		}

		if entry.IsDir() {
			err = processDirectory(baseDir, filepath.Join(currentDir, entry.Name()), ignoreList, out, depth+1)
			if err != nil {
				return err
			}
		} else {
			err = processFile(baseDir, filepath.Join(currentDir, entry.Name()), out)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processFile(baseDir, filePath string, out *os.File) error {
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
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
