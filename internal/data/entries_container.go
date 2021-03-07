package data

import "github.com/pkg/errors"

type EntriesContainer struct {
	dataProvider                     Provider
	entries                          map[EntryType][]Entry
	changeListenersCallbackFunctions []func()
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
		return errors.New("Entry type with name '" + entryTypeToAdd.Name + "' already exists")
	}
	container.entries[entryTypeToAdd] = []Entry{}
	container.notifyListenersAboutChange()
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

func (container *EntriesContainer) notifyListenersAboutChange() {
	for _, callback := range container.changeListenersCallbackFunctions {
		callback()
	}
}

func (container *EntriesContainer) DeleteEntryType(typeName string) error {
	if container.typeWithNameExists(typeName) {
		container.deleteEntryTypeWithName(typeName)
		container.notifyListenersAboutChange()
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

func (container *EntriesContainer) UpdateEntryType(nameOfTypeToUpdate string, typeToReplaceWith EntryType) error {
	if typeToReplaceWith.Name == "" {
		return errors.New("Cannot update entry type with name '" + nameOfTypeToUpdate + "' to type with an empty name")
	}
	return container.tryUpdatingEntryType(nameOfTypeToUpdate, typeToReplaceWith)
}

func (container *EntriesContainer) tryUpdatingEntryType(nameOfTypeToUpdate string, typeToReplaceWith EntryType) error {
	for entryType, entries := range container.entries {
		if entryType.Name == nameOfTypeToUpdate {
			delete(container.entries, entryType)
			container.entries[typeToReplaceWith] = entries
			container.notifyListenersAboutChange()
			return nil
		}
	}
	return errors.New("Cannot update entry type '" + nameOfTypeToUpdate + "' as no such type exists")
}

func (container *EntriesContainer) EntryTypeWithName(typeName string) (EntryType, error) {
	for entryType, _ := range container.entries {
		if entryType.Name == typeName {
			return entryType, nil
		}
	}
	return EntryType{}, errors.New("Cannot retrieve entry type with name '" + typeName + "' as such entry type doesn't exist")
}

func (container *EntriesContainer) EntriesGroupedByType() map[EntryType][]Entry {
	entriesToReturn := make(map[EntryType][]Entry, len(container.entries))
	for entryType, entries := range container.entries {
		entriesToReturn[entryType] = entries
	}
	return entriesToReturn
}

func (container *EntriesContainer) SubscribeToChanges(callbackFunction func()) {
	container.changeListenersCallbackFunctions = append(container.changeListenersCallbackFunctions, callbackFunction)
}

func (container *EntriesContainer) AmountOfTypes() int {
	return len(container.entries)
}
