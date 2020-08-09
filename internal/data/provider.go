package data

type Provider interface {
	SaveEntries(map[EntryType][]Entry) error
	LoadEntries() (map[EntryType][]Entry, error)
}
