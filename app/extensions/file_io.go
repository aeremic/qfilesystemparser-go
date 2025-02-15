package extensions

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func ReadInputAsString(reader *bufio.Reader) string {
	if text, error := reader.ReadString('\n'); error != nil {
		fmt.Println("Error reading input: ", error)
		return ""
	} else {
		return strings.Replace(text, "\n", "", -1)
	}
}

func ReadInputAsInt(reader *bufio.Reader) int {
	if value, error :=
		strconv.Atoi(ReadInputAsString(reader)); error != nil {
		fmt.Println("Error reading input: ", error)
		return 0
	} else {
		return value
	}
}

func WriteString(writer *bufio.Writer, output string) {
	_, writeError :=
		writer.WriteString(output)
	if writeError != nil {
		panic(writeError)
	}
}
