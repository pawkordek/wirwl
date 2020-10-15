package data

type EntriesContainer struct {
	dataProvider Provider
	entries      map[EntryType][]Entry
}

func NewEntriesContainer(dataProvider Provider) *EntriesContainer {
	return &EntriesContainer{dataProvider: dataProvider}
}

func (controller *EntriesContainer) LoadData() error {
	entries, err := controller.dataProvider.LoadEntries()
	controller.entries = entries
	return err
}

func (controller *EntriesContainer) SaveData() error {
	err := controller.dataProvider.SaveEntries(controller.entries)
	return err
}
