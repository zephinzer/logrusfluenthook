package integrationtest

import (
	"path"
	"runtime"
	"strings"
	"github.com/phayes/freeport"
)

func getCurrentFileDirectory() string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename)
	}
	return ""
}

func getFreePort() int {
	freePort, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	return freePort
}

// getLastLineOfString returns the number of lines and
// the last line of the :input string, meant for retrieving
// the latest line of logs from the fluentd service
func getLastLineOfString(input string) (uint, string) {
	normalisedFormat :=
		strings.Trim(
			strings.Replace(input, "\r\n", "\n", -1),
			" \n",
		)
	splitInputByLine :=
		strings.Split(normalisedFormat, "\n")
	numberOfLines := len(splitInputByLine)
	lastLine := splitInputByLine[numberOfLines - 1]
	return uint(numberOfLines), lastLine
}