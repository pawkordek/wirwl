package data

import "github.com/pkg/errors"

type EntriesContainer struct {
	dataProvider Provider
	entries      map[EntryType][]Entry
}

func NewEntriesContainer(dataProvider Provider) *EntriesContainer {
	return &EntriesContainer{entries: map[EntryType][]Entry{}, dataProvider: dataProvider}
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

func (container *EntriesContainer) AddEntryType(entryTypeToAdd EntryType) error {
	if entryTypeToAdd.Name == "" {
		return errors.New("Cannot add entry type with an empty name")
	} else if container.typeWithNameExists(entryTypeToAdd.Name) {
		return errors.New("Entry type with name " + entryTypeToAdd.Name + " already exists")
	}
	container.entries[entryTypeToAdd] = []Entry{}
	return nil
}

func (container *EntriesContainer) typeWithNameExists(nameToCheck string) bool {
	for entryType, _ := range container.entries {
		if entryType.Name == nameToCheck {
			return true
		}
	}
	return false
}

func (container *EntriesContainer) DeleteEntryType(typeName string) error {
	if container.typeWithNameExists(typeName) {
		container.deleteEntryTypeWithName(typeName)
		return nil
	}
	return errors.New("Cannot delete an entry type with name '" + typeName + "' as there is no such type")
}

func (container *EntriesContainer) deleteEntryTypeWithName(typeName string) {
	for entryType, _ := range container.entries {
		if entryType.Name == typeName {
			delete(container.entries, entryType)
		}
	}
}
