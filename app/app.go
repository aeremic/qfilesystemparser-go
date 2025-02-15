package main

import (
	"aeremic/qfilesystemparser/logic"
	"context"
	"fmt"
	"log"
	"net/rpc"
)

// App struct
type App struct {
	ctx    context.Context
	client *rpc.Client
	error  error
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	client, error := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if error != nil {
		log.Fatal("Error on dialing: ", error)
	}

	a.client = client
	a.error = error
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ShowHelloWorld(name string) string {
	return fmt.Sprintf("Hello world to %s", name)
}

func (a *App) StartParsing(filesPath string,
	checkInterval int,
	maximumNumberOfProcessingJobs int,
	maximumExecutionCount int,
) string {
	if filesPath == "" ||
		checkInterval == 0 ||
		maximumNumberOfProcessingJobs == 0 ||
		maximumExecutionCount == 0 {
		return "Invalid input."
	}

	inputDataArgs := &logic.InputDataArgs{
		FilesPath:                     filesPath,
		CheckInterval:                 checkInterval,
		MaximumNumberOfProcessingJobs: maximumNumberOfProcessingJobs,
		MaximumExecutedCount:          maximumExecutionCount}

	setInputDataReply := &logic.MethodCallResult{}
	if a.client.Call("Parser.SetInputData", inputDataArgs, &setInputDataReply); a.error != nil {
		log.Fatal("Error on client call: ", a.error)
	}

	doParsingReply := &logic.MethodCallResult{}
	parsingArgs := &logic.ParsingArgs{}
	if a.client.Call("Parser.DoParsing", parsingArgs, &doParsingReply); a.error != nil {
		log.Fatal("Error on client call: ", a.error)
	}

	if doParsingReply.IsSuccess {
		return "Parsing successfull."
	} else {
		return "Parsing failed."
	}
}
