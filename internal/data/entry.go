package data

type Entry struct {
	Id                         int
	Status                     string
	Title                      string
	Completion                 int
	AmountOfElementsToComplete int
	Score                      int
	Link                       string
	Description                string
	Comment                    string
	MediaType                  string
	Tags                       string
}

type EntryType struct {
	Name       string
	ImageQuery string
}
