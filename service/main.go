package main

import (
	"aeremic/qfilesystemparser/extensions"
	"aeremic/qfilesystemparser/logic"
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {
	rpcParser := new(logic.Parser)
	rpc.Register(rpcParser)
	rpc.HandleHTTP()
	listen, error := net.Listen("tcp", ":1234")
	if error != nil {
		log.Fatal("Listen error: ", error)
	}

	go http.Serve(listen, nil)

	reader := bufio.NewReader(os.Stdin)
	extensions.ReadInputAsString(reader)

	fmt.Println("Ending...")
}
