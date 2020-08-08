package data

type Provider interface {
	SaveEntriesToDb(table string, entries []Entry) error
	LoadEntriesFromDb(table string) ([]Entry, error)
	SaveEntriesTypesToDb(entriesTypes []EntryType) error
	LoadEntriesTypesFromDb() ([]EntryType, error)
	SaveEntries(map[EntryType][]Entry) error
}
