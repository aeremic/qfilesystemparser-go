package logic

import (
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

func (wgw *WaitGroupWrapper) Count() int {
	return int(atomic.LoadInt64(&wgw.count))
}

var waitGroup sync.WaitGroup
var mutex sync.Mutex

type Parser int

var FilesPath string
var MaximumNumberOfProcessingJobs int
var CheckInterval int
var MaximumExecutedCount int

type InputDataArgs struct {
	FilesPath                     string
	MaximumNumberOfProcessingJobs int
	CheckInterval                 int
	MaximumExecutedCount          int
}

type ParsingArgs struct{}

func (t *Parser) SetInputDataViaRpc(args *InputDataArgs, reply *int) error {
	FilesPath = args.FilesPath
	MaximumNumberOfProcessingJobs = args.MaximumNumberOfProcessingJobs
	CheckInterval = args.CheckInterval
	MaximumExecutedCount = args.MaximumExecutedCount

	*reply = args.MaximumNumberOfProcessingJobs + args.MaximumExecutedCount

	return nil
}

func (t *Parser) DoParsing(args *ParsingArgs, reply *int) error {
	executedCounter := 0
	for executedCounter < MaximumExecutedCount {
		waitGroupWrapper := WaitGroupWrapper{&waitGroup, int64(0)}
		waitGroupWrapper.Add(1)
		walkDirectory(&waitGroupWrapper, MaximumNumberOfProcessingJobs, FilesPath)
		waitGroupWrapper.Wait()

		time.Sleep(time.Duration(CheckInterval) * time.Second)

		executedCounter++
	}

	*reply = 1

	return nil
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

	var quest Quest

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

			mutex.Lock()
			if int64(waitGroupWrapper.Count()) < int64(maximumNumberOfProcessingJob) {
				go walkDirectory(waitGroupWrapper, maximumNumberOfProcessingJob, path)
			} else {
				walkDirectory(waitGroupWrapper, maximumNumberOfProcessingJob, path)
			}
			mutex.Unlock()

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
