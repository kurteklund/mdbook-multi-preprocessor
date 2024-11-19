package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		log.Printf("Error: (%v)\n", e)
		panic(e)
	}
}

func CreateTextFile(filetPath string, content string) {
	f, err := os.Create(filetPath)
	check(err)

	n, err := f.WriteString(content + "\n")
	check(err)
	if n < 0 {
		panic("tRAMS")
	}
	f.Sync()
}

func ReadTextFile(path string) string {
	result := ""
	f, _ := os.Open(path)
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		result += scanner.Text()
	}

	return result
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

func GetConditionalRegions(topItem MdBookTopItem) []string {
	var preprocessor = topItem.Config.Preprocessor
	if preprocessor.Test != nil && preprocessor.Test.ConditionalRegions != nil {
		return preprocessor.Test.ConditionalRegions
	}

	return []string{}
}

func main() {
	const logFileName = "/tmp/multi/log.txt"
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	if len(os.Args) > 1 {
		if os.Args[1] == "supports" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	jsonText := ReadTextFile("/tmp/multi/input.json")
	// jsonText := readJsonFromStdIn()
	//	CreateTextFile("/tmp/multi/input.json", jsonText)
	log.Println("After readJsonFromStdIn()")

	var book []MdBookTopItem
	errJson := json.Unmarshal([]byte(jsonText), &book)
	check(errJson)
	if len(book) == 2 {
		// The input json is an slice with 2 items.
		// The first item is configuration, parameters to the preprocessor and other stuff
		// The second item is the "content" of the book, the part that should be exported
		conditionalRegions := GetConditionalRegions(book[0])
		bookSections := book[1]
		processSections(&bookSections, conditionalRegions)
		// writeBookSectionsToFile(bookSections, "/tmp/mdbook_out.json")
		writeBookSectionsStdOut(bookSections)
	}
}
