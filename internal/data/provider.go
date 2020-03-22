package data

type Provider interface {
	SaveEntriesToDb(table string, entries []Entry) error
	LoadEntriesFromDb(table string) ([]Entry, error)
	SaveEntriesTypesToDb(entriesTypes []string) error
	LoadEntriesTypesFromDb() ([]string, error)
}
