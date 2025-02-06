package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
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
