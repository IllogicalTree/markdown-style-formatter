package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type change struct {
	filepath       string
	markdown       string
	match          string
	convertedMatch string
}

func findFilesWithExtension(rootDirectory, extension string) []string {
	var files []string
	filepath.WalkDir(rootDirectory, func(s string, directory fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(directory.Name()) == extension {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func convertToTitleCase(s string) string {
	if strings.Contains(s, "http") {
		return s
	}
	return strings.Title(s)
}

func convertToSentenceCase(s string) string {
	s = strings.ToLower(s)
	r := []rune(s)
	i := strings.Index(s, "# ")
	r[i+2] = unicode.ToUpper(r[i+2])
	return string(r)
}

func stripTrailingPeriod(s string) string {
	return s[:len(s)-1]
}

func main() {

	changes := []change{}

	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run . <filePath>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	for _, file := range findFilesWithExtension(filePath, ".mdx") {

		fileContents, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		markdown := string(fileContents)

		makeConversions := func(markdown string, filepath string, r *regexp.Regexp, convert func(string) string) {
			matches := r.FindAllString(markdown, -1)
			for _, match := range matches {
				match = strings.TrimSpace(match)
				convertedMatch := convert(match)

				// Headings should not end with a period.
				convertedMatch = strings.TrimSuffix(convertedMatch, ".")

				if match != convertedMatch {
					changes = append(changes, change{filepath, markdown, match, convertedMatch})
				}
			}
		}

		// Headings should use Title Capitalization Like This.
		r := regexp.MustCompile(`\n(#{1}\s)(.*)`)
		makeConversions(markdown, file, r, convertToTitleCase)

		// Subheadings, anything less than h1 or markdown level 1: #, should use Sentence capitalization like this.
		r = regexp.MustCompile(`\n(#{2-6}\s)(.*)`)
		makeConversions(markdown, file, r, convertToSentenceCase)
	}

	for _, change := range changes {
		markdown := strings.ReplaceAll(change.markdown, change.match, change.convertedMatch)
		newContents := []byte(markdown)
		err := os.WriteFile(change.filepath, newContents, 0644)
		fmt.Println("Wrote to file: " + change.filepath + " : " + change.match + " -> " + change.convertedMatch)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Done!")
}
