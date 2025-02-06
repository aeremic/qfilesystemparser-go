package logic

import (
	"time"
)

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

func (t *Parser) DoParsingViaRpc(args *ParsingArgs, reply *int) error {
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
