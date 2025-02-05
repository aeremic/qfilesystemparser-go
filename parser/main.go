package main

import (
	"aeremic/qfilesystemparser/common"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type WaitGroupWrapper struct {
	*sync.WaitGroup
	count int64
}

func (wgw *WaitGroupWrapper) Add(delta int) {
	atomic.AddInt64(&wgw.count, int64(delta))
	wgw.WaitGroup.Add(delta)
}

func (wgw *WaitGroupWrapper) Done() {
	atomic.AddInt64(&wgw.count, -1)
	wgw.WaitGroup.Done()
}

func (wgw *WaitGroupWrapper) Wait() {
	wgw.WaitGroup.Wait()
}

func (wgw *WaitGroupWrapper) GetCount() int {
	return int(atomic.LoadInt64(&wgw.count))
}

var waitGroup sync.WaitGroup

// Reads data required for input.
func readRequiredInputData() (string, int, int, int) {
	// reader := bufio.NewReader(os.Stdin)

	fmt.Println("------------------------------------------------------------")

	fmt.Print("Enter files path for parsing: ")
	// filespath := extensions.ReadInputAsString(reader)
	filespath := "/Users/eremic/qfilesystemparser-go/Files"
	fmt.Println("Files path: ", filespath)

	fmt.Print("Enter check interval: ")
	// checkInterval := extensions.ReadInputAsInt(reader)
	checkInterval := 5
	fmt.Println("Check interval: ", checkInterval)

	fmt.Print("Enter maximum number of processing jobs: ")
	// maximumNumberOfProcessingJobs := extensions.ReadInputAsInt(reader)
	maximumNumberOfProcessingJobs := 1
	fmt.Println("Maximum number of processing jobs: ", maximumNumberOfProcessingJobs)

	fmt.Print("Enter maximum number of executions: ")
	// maximumNumberOfProcessingJobs := extensions.ReadInputAsInt(reader)
	maximumExecutedCount := 3
	fmt.Println("Maximum number of executions: ", maximumExecutedCount)

	fmt.Println("------------------------------------------------------------")

	return filespath, checkInterval, maximumNumberOfProcessingJobs, maximumExecutedCount
}

func parseJsonFile(path string) {
	jsonFile, error := os.Open(path)
	defer jsonFile.Close()

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

	fmt.Println("Number of components in ", path, " file: ", len(quest.Components))
}

func walkDirectory(waitGroupWrapper *WaitGroupWrapper,
	maximumNumberOfProcessingJob int, filesPath string) {
	defer waitGroupWrapper.Done()

	visit := func(path string, fileInfo os.FileInfo, error error) error {
		if fileInfo.IsDir() && path != filesPath {
			waitGroupWrapper.Add(1)

			if waitGroupWrapper.count < int64(maximumNumberOfProcessingJob) {
				go walkDirectory(waitGroupWrapper, maximumNumberOfProcessingJob, path)
			} else {
				walkDirectory(waitGroupWrapper, maximumNumberOfProcessingJob, path)
			}

			return filepath.SkipDir
		}

		if fileInfo.Mode().IsRegular() {
			if filepath.Ext(path) == ".json" {
				parseJsonFile(path)
			}
		}

		return nil
	}

	filepath.Walk(filesPath, visit)
}

func main() {
	filesPath, maximumNumberOfProcessingJob, checkInterval, maximumExecutedCount := readRequiredInputData()

	executedCounter := 0
	for executedCounter < maximumExecutedCount {
		waitGroupWrapper := WaitGroupWrapper{&waitGroup, int64(0)}
		waitGroupWrapper.Add(1)
		walkDirectory(&waitGroupWrapper, maximumNumberOfProcessingJob, filesPath)
		waitGroupWrapper.Wait()

		time.Sleep(time.Duration(checkInterval) * time.Second)

		executedCounter++
	}
}
