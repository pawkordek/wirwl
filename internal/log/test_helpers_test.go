package log

import (
	"os"
)

//Since it's not possible to compare the text written to standard output to check whether logger works correctly,
//the general idea is to mock this output and store whatever was written to it into a field in a struct using a callback
//function, which can be then compared to the expected value

type writeResult struct {
	value string
}

func newMockWrite() *writeResult {
	return &writeResult{""}
}

func (write *writeResult) storeWriteValue(s string) {
	write.value = s
}

func (write writeResult) String() string {
	return write.value
}

type mockOutput struct {
	writeCallback func(string)
}

func (output mockOutput) Write(p []byte) (n int, err error) {
	output.writeCallback(string(p))
	return len(p), nil
}

func setupLoggerTesting() (*writeResult, func()) {
	writeResult := newMockWrite()
	output := mockOutput{writeResult.storeWriteValue}
	SetOutput(output)
	return writeResult, cleanupLoggerTesting
}

func cleanupLoggerTesting() {
	SetOutput(os.Stdout)
}
