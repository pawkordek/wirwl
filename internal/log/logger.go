package log

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Info(info string) {
	log.Print("INFO: " + info)
}

func Error(err error) {
	printError("ERROR:", err)
}

func printError(textBefore string, err error) {
	printableError := getPrintableError(err)
	log.Printf(textBefore+" %+v", printableError)
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
