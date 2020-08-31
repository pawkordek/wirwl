package data

import "fmt"

type EntryStatus string

const (
	InProgressStatus EntryStatus = "In progress"
	CompletedStatus  EntryStatus = "Completed"
	OnHoldStatus     EntryStatus = "On hold"
	DroppedStatus    EntryStatus = "Dropped"
	PlannedStatus    EntryStatus = "Planned"
)

type Entry struct {
	Id                              int
	Status                          EntryStatus
	Title                           string
	ElementsCompleted               int
	TotalAmountOfElementsToComplete int
	Score                           int
	Link                            string
	Description                     string
	Comment                         string
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
