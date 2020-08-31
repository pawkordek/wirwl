package data

import "fmt"

type Entry struct {
	Id                              int
	Status                          string
	Title                           string
	ElementsCompleted               int
	TotalAmountOfElementsToComplete int
	Score                           int
	Link                            string
	Description                     string
	Comment                         string
	MediaType                       string
	Tags                            string
	ImageQuery                      string
}

func (entry Entry) String() string {
	return fmt.Sprintf("%#v", entry)
}

type EntryType struct {
	Name       string
	ImageQuery string
}

func (entryType EntryType) String() string {
	return fmt.Sprintf("%#v", entryType)
}
