package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// The input json is an array with 2 items.
// The first item is configuration, parameters to the preprocessor and other stuff
// The second item is the "content" of the book, the part that should be exported
// This is strange, why in same array?
type MdBookTopItem struct {
	Root          string          `json:"root,omitempty"`
	Config        *MdBookConfig   `json:"config,omitempty"`
	Renderer      string          `json:"renderer,omitempty"`
	MdBookVersion string          `json:"mdbook_version,omitempty"`
	Sections      []MdBookSection `json:"sections,omitempty"`
	NonExhaustive *string         `json:"__non_exhaustive"`
}

type MdBookConfig struct {
	Book         MdBookConfigBook         `json:"book"`
	Build        MdBookConfigBuild        `json:"build"`
	Output       MdBookConfigOutput       `json:"output"`
	Preprocessor MdBookConfigPreprocessor `json:"preprocessor"`
}

type MdBookConfigBook struct {
	Authors      []string `json:"authors"`
	Language     string   `json:"language"`
	Multilingual bool     `json:"multilingual"`
	Src          string   `json:"src"`
	Title        string   `json:"title"`
}

type MdBookConfigBuild struct {
	BuildDir                string   `json:"build-dir"`
	CreateMissing           bool     `json:"create-missing"`
	ExtraWatchDirs          []string `json:"extra-watch-dirs"`
	UseDefaultPreprocessors bool     `json:"use-default-preprocessors"`
}

type MdBookConfigOutput struct {
	Html MdBookConfigOutputHtml `json:"html"`
}

type MdBookConfigOutputHtml struct {
	AdditionalCss []string `json:"additional-css"`
}

type MdBookConfigPreprocessor struct {
	Example *MdBookConfigPreprocessorX `json:"example,omitempty"`
	Trams   *MdBookConfigPreprocessorX `json:"trams,omitempty"`
}

type MdBookConfigPreprocessorX struct {
	Before   []string `json:"before,omitempty"`
	After    []string `json:"after,omitempty"`
	Command  string   `json:"command"`
	Renderer []string `json:"renderer,omitempty"`
}

type MdBookSection struct {
	Chapter MdBookChapter `json:"Chapter"`
}

type MdBookChapter struct {
	Name        string          `json:"name"`
	Content     string          `json:"content"`
	Number      []int           `json:"number"`
	SubItems    []MdBookSection `json:"sub_items"`
	Path        string          `json:"path"`
	SourcePath  string          `json:"source_path"`
	ParentNames []string        `json:"parent_names"`
}

func check(e error) {
	if e != nil {
		fmt.Printf("Error: (%v)\n", e)
		panic(e)
	}
}

func readFileLinesFile(filePath string) []string {
	readFile, err := os.Open(filePath)

	check(err)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func readTextFile(filePath string) string {
	lines := readFileLinesFile(filePath)
	text := ""

	for _, line := range lines {
		text += line
	}

	return text
}

func writeBookSectionsToFile(bookSections MdBookTopItem, filename string) {
	jsonData, err := json.MarshalIndent(bookSections, "", "    ")
	check(err)
	f, err := os.Create(filename)
	check(err)
	f.Write(jsonData)
}

func writeBookSectionsStdOut(bookSections MdBookTopItem) {
	jsonData, err := json.Marshal(bookSections)
	check(err)
	fmt.Println(string(jsonData))
}

func readJsonFromStdIn() string {
	f, err := os.Create("/tmp/kurt_in.json")
	check(err)
	tempWriter := bufio.NewWriter(f)

	scanner := bufio.NewScanner(os.Stdin)
	jsonText := ""
	for scanner.Scan() {
		jsonText += scanner.Text() + "\n"
		_, err := tempWriter.WriteString(scanner.Text() + "\n")
		check(err)
	}

	tempWriter.Flush()
	return jsonText
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "supports" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	//jsonText := readTextFile("/tmp/sample.json")
	jsonText := readJsonFromStdIn()

	var book []MdBookTopItem
	errJson := json.Unmarshal([]byte(jsonText), &book)
	check(errJson)
	if len(book) == 2 {
		// The input json is an slice with 2 items.
		// The first item is configuration, parameters to the preprocessor and other stuff
		// The second item is the "content" of the book, the part that should be exported
		bookSections := book[1]
		processSections(&bookSections)
		// writeBookSectionsToFile(bookSections, "/tmp/mdbook_out.json")
		writeBookSectionsStdOut(bookSections)
	}
}
