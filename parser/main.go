package main

import (
	"aeremic/qfilesystemparser/common"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func readRequiredInputData() (string, int, int) {
	// reader := bufio.NewReader(os.Stdin)

	fmt.Println("------------------------------------------------------------")

	fmt.Print("Enter files path for parsing: ")
	// filespath := extensions.ReadInputAsString(reader)
	filespath := "/Users/eremic/qfilesystemparser-go/Files"
	fmt.Println("Files path: ", filespath)

	fmt.Print("Enter check interval: ")
	// checkInterval := extensions.ReadInputAsInt(reader)
	checkInterval := 1000
	fmt.Println("Check interval: ", checkInterval)

	fmt.Print("Enter maximum number of processing jobs: ")
	// maximumNumberOfProcessingJobs := extensions.ReadInputAsInt(reader)
	maximumNumberOfProcessingJobs := 1
	fmt.Println("Maximum number of processing jobs: ", maximumNumberOfProcessingJobs)

	fmt.Println("------------------------------------------------------------")

	return filespath, checkInterval, maximumNumberOfProcessingJobs
}

func main() {
	filesPath, _, _ := readRequiredInputData()

	fileSystem := os.DirFS(filesPath)
	fs.WalkDir(fileSystem, ".", func(path string, directory fs.DirEntry, error error) error {
		if filepath.Ext(path) == ".json" {
			jsonFile, error := os.Open(filesPath + "/" + path)
			if error != nil {
				fmt.Println("Error on opening json file at path: ", path,
					" Error details: ", error)
			}

			byteArray, error := io.ReadAll(jsonFile)
			if error != nil {
				fmt.Println("Error on reading json file: ", path,
					" Error details: ", error)
			}

			var quest common.Quest

			var errorOnJsonParsing = json.Unmarshal(byteArray, &quest)
			if errorOnJsonParsing != nil {
				fmt.Println("Error on parsing json file: ", path,
					" Error details: ", errorOnJsonParsing)
			}
			defer jsonFile.Close()

			fmt.Println("Number of components in ", path, " file: ", len(quest.Components))
		}

		return nil
	})
}
