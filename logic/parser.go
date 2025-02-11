package logic

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Parser int

var FilesPath string
var MaximumNumberOfProcessingJobs int
var CheckInterval int
var MaximumExecutedCount int

type MethodCallResult struct {
	IsSuccess    bool
	ErrorMessage error
}

type InputDataArgs struct {
	FilesPath                     string
	MaximumNumberOfProcessingJobs int
	CheckInterval                 int
	MaximumExecutedCount          int
}

type ParsingArgs struct{}

func (t *Parser) SetInputData(args *InputDataArgs, reply *MethodCallResult) error {
	FilesPath = args.FilesPath
	MaximumNumberOfProcessingJobs = args.MaximumNumberOfProcessingJobs
	CheckInterval = args.CheckInterval
	MaximumExecutedCount = args.MaximumExecutedCount

	*reply = MethodCallResult{true, nil}

	return nil
}

var quitChannel = make(chan bool)

func (t *Parser) DoParsing(args *ParsingArgs, reply *MethodCallResult) error {
	errorChannel := make(chan error)

	file, error := os.Create("output.txt")
	if error != nil {
		panic(error)
	}
	defer file.Close()

	go func() {
		writer := bufio.NewWriter(file)
		executedCounter := 0
		for executedCounter < MaximumExecutedCount {
			select {
			case <-quitChannel:
				fmt.Println("Parsing stopped.")
				return
			default:
				waitGroupWrapper := WaitGroupWrapper{&waitGroup, int64(0)}
				waitGroupWrapper.Add(1)
				walkDirectory(writer, &waitGroupWrapper, MaximumNumberOfProcessingJobs, FilesPath)
				waitGroupWrapper.Wait()

				time.Sleep(time.Duration(CheckInterval) * time.Second)

				executedCounter++
			}
		}

		writer.Flush()

		*reply = MethodCallResult{true, nil}
		errorChannel <- nil
	}()

	return <-errorChannel
}

type StopParsingArgs struct{}

func (t *Parser) StopParsing(arg StopParsingArgs, reply *MethodCallResult) error {
	quitChannel <- true

	*reply = MethodCallResult{true, nil}

	return nil
}
