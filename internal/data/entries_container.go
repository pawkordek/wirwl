package data

type EntriesContainer struct {
	dataProvider Provider
	entries      map[EntryType][]Entry
}

func NewEntriesContainer(dataProvider Provider) *EntriesContainer {
	return &EntriesContainer{dataProvider: dataProvider}
}

func (container *EntriesContainer) LoadData() error {
	entries, err := container.dataProvider.LoadEntries()
	container.entries = entries
	return err
}

func (container *EntriesContainer) SaveData() error {
	err := container.dataProvider.SaveEntries(container.entries)
	return err
}
