package logic

import (
	"fmt"
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
	go func() {
		executedCounter := 0
		for executedCounter < MaximumExecutedCount {
			select {
			case <-quitChannel:
				fmt.Println("Parsing stopped.")
				return
			default:
				waitGroupWrapper := WaitGroupWrapper{&waitGroup, int64(0)}
				waitGroupWrapper.Add(1)
				walkDirectory(&waitGroupWrapper, MaximumNumberOfProcessingJobs, FilesPath)
				waitGroupWrapper.Wait()

				time.Sleep(time.Duration(CheckInterval) * time.Second)

				executedCounter++
			}
		}
	}()

	*reply = MethodCallResult{true, nil}

	return nil
}

type StopParsingArgs struct{}

func (t *Parser) StopParsing(arg StopParsingArgs, reply *MethodCallResult) error {
	quitChannel <- true

	*reply = MethodCallResult{true, nil}

	return nil
}
