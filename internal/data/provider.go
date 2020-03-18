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
}

type Provider interface {
	SaveEntriesToDb(table string, entries []Entry) error
	LoadEntriesFromDb(table string) ([]Entry, error)
	SaveEntriesTypesToDb(entriesTypes []string) error
	LoadEntriesTypesFromDb() ([]string, error)
}
