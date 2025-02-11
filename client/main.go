package main

import (
	"aeremic/qfilesystemparser/extensions"
	"aeremic/qfilesystemparser/logic"
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

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
	// maximumExecutedCount := extensions.ReadInputAsInt(reader)
	maximumExecutedCount := 5
	fmt.Println("Maximum number of executions: ", maximumExecutedCount)

	fmt.Println("------------------------------------------------------------")

	return filespath, checkInterval, maximumNumberOfProcessingJobs, maximumExecutedCount
}

func main() {
	filesPath, maximumNumberOfProcessingJob,
		checkInterval, maximumExecutedCount := readRequiredInputData()

	client, error := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if error != nil {
		log.Fatal("Error on dialing: ", error)
	}

	inputDataArgs := &logic.InputDataArgs{
		FilesPath:                     filesPath,
		CheckInterval:                 checkInterval,
		MaximumNumberOfProcessingJobs: maximumNumberOfProcessingJob,
		MaximumExecutedCount:          maximumExecutedCount}

	setInputDataReply := &logic.MethodCallResult{}
	if client.Call("Parser.SetInputData", inputDataArgs, &setInputDataReply); error != nil {
		log.Fatal("Error on client call: ", error)
	}

	doParsingReply := &logic.MethodCallResult{}
	parsingArgs := &logic.ParsingArgs{}
	if client.Call("Parser.DoParsing", parsingArgs, &doParsingReply); error != nil {
		log.Fatal("Error on client call: ", error)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter STOP to cancel parsing: ")
	stopInput := extensions.ReadInputAsString(reader)
	if stopInput == "STOP" {
		stopParsingArgs := &logic.StopParsingArgs{}
		stopParsingReply := &logic.MethodCallResult{}
		if client.Call("Parser.StopParsing", stopParsingArgs, &stopParsingReply); error != nil {
			log.Fatal("Erro on stoping client call: ", error)
		}
	}
}
