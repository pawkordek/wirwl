package log

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Info(info string) {
	log.Print("INFO: " + info)
}

func Error(err error) {
	printError("ERROR: ", err)
}

//Should only be used in tests to quickly fail and report an error.
//In application code all errors should be handled if possible, otherwise they should bubble up to the main function,
//where application can exit.
func Fatal(err error) {
	printError("FATAL: ", err)
	os.Exit(1)
}

//Should only be used in rare cases when it makes sense to panic on a programmer error so that the cause can be logged
func Panic(err error) {
	printError("PANIC: ", err)
	panic(err)
}

func printError(textBefore string, err error) {
	printableError := getPrintableError(err)
	log.Printf(textBefore+"%+v", printableError)
}

func getPrintableError(err error) string {
	if stackTrace, ok := err.(stackTracer); ok {
		buffer := bytes.Buffer{}
		for _, frame := range stackTrace.StackTrace() {
			buffer.WriteString(fmt.Sprintf("\n%+v", frame))
		}
		return err.Error() + buffer.String()
	} else {
		return err.Error()
	}
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}
