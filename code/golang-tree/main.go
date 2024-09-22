package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	showAllFiles    = flag.Bool("a", false, "Show all files and directories, including hidden ones")
	showDirectories = flag.Bool("d", false, "Only list directories")
	maxDepth        = flag.Int("L", -1, "Limit the level of recursion")
	showFullPath    = flag.Bool("f", false, "Show full path for each file or directory")
	pattern         = flag.String("P", "", "Show only files/directories that match the provided pattern")
	excludePattern  = flag.String("I", "", "Exclude files/directories that match the provided pattern")
	outputFile      = flag.String("o", "", "Redirect output to the specified file")
	noIndent        = flag.Bool("i", false, "Do not display indentation lines")
)

func main() {
	flag.Parse()

	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	var output *os.File
	var err error
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	err = tree(dir, "", 0, output)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func tree(path string, prefix string, level int, output *os.File) error {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if !*showAllFiles && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if *showDirectories && !entry.IsDir() {
			continue
		}
		if *excludePattern != "" {
			excludeRegexp, err := regexp.Compile(*excludePattern)
			if err != nil {
				return err
			}
			if excludeRegexp.MatchString(entry.Name()) {
				continue
			}
		}
		if *pattern != "" {
			patternRegexp, err := regexp.Compile(*pattern)
			if err != nil {
				return err
			}
			if !patternRegexp.MatchString(entry.Name()) {
				continue
			}
		}

		// Print the current entry
		if *showFullPath {
			fmt.Fprintln(output, getPrefix(prefix, i, len(entries))+filepath.Join(path, entry.Name()))
		} else {
			fmt.Fprintln(output, getPrefix(prefix, i, len(entries))+entry.Name())
		}

		if entry.IsDir() {
			if *maxDepth < 0 || level < *maxDepth {
				newPrefix := getNextPrefix(prefix, i, len(entries))
				err := tree(filepath.Join(path, entry.Name()), newPrefix, level+1, output)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getPrefix(prefix string, index, total int) string {
	if *noIndent {
		return ""
	}
	if index == total-1 {
		return prefix + "└── "
	}
	return prefix + "├── "
}

func getNextPrefix(prefix string, index, total int) string {
	if *noIndent {
		return ""
	}
	if index == total-1 {
		return prefix + "    "
	}
	return prefix + "│   "
}
