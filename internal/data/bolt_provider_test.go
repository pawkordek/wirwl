package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wirwl/internal/log"
)

func TestDbOperationsOnEntries(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	entriesToSave := GetTestEntries()
	dataProvider := NewBoltProvider(testDbPath)
	err := dataProvider.SaveEntries(entriesToSave)
	if err != nil {
		log.Fatal(err)
	}
	loadedEntries, err := dataProvider.LoadEntries()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, entriesToSave[comicsEntryType], loadedEntries[comicsEntryType])
	assert.Equal(t, entriesToSave[musicEntryType], loadedEntries[musicEntryType])
	assert.Equal(t, entriesToSave[videoEntryType], loadedEntries[videoEntryType])
	DeleteTestDb()
}

func TestThatTryingToLoadDataIntoNonExistingTableReturnsError(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	types, err := dataProvider.LoadEntriesFromDb("non existing table")
	assert.Nil(t, types)
	assert.Contains(t, err.Error(), "An error occurred when loading entries from table with name non existing table. No such table")
	DeleteTestDb()
}

func TestDbOperationsOnEntriesTypes(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	entriesTypes := GetEntriesTypes()
	dataProvider := NewBoltProvider(testDbPath)
	err := dataProvider.SaveEntriesTypesToDb(entriesTypes)
	if err != nil {
		log.Fatal(err)
	}
	types, err := dataProvider.LoadEntriesTypesFromDb()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, types[0], entriesTypes[0])
	assert.Equal(t, types[1], entriesTypes[1])
	assert.Equal(t, types[2], entriesTypes[2])
	DeleteTestDb()
}

func TestThatTryingToLoadEntriesFromEmptyDbReturnsEmptySlice(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	entriesToSave := GetTestEntriesToSave()
	err := dataProvider.SaveEntriesToDb("entries", entriesToSave)
	if err != nil {
		log.Fatal(err)
	}
	types, err := dataProvider.LoadEntriesTypesFromDb()
	assert.Equal(t, 0, len(types))
	assert.Nil(t, err)
	DeleteTestDb()
}

func TestThatSavingEmptyEntriesSliceCreatesTable(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	err := dataProvider.SaveEntriesToDb("new table", []Entry{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = dataProvider.LoadEntriesFromDb("new table")
	assert.Nil(t, err)
	DeleteTestDb()
}

func TestThatWhenSavingEntriesPreviousDataInDbIsRemoved(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	entriesToSave := GetTestEntriesToSave()
	err := dataProvider.SaveEntriesToDb("entries", entriesToSave)
	if err != nil {
		log.Fatal(err)
	}
	err = dataProvider.SaveEntriesToDb("entries", []Entry{})
	if err != nil {
		log.Fatal(err)
	}
	loadedEntries, err := dataProvider.LoadEntriesFromDb("entries")
	assert.Empty(t, loadedEntries)
	DeleteTestDb()
}

func TestThatWhenSavingEntriesTypesPreviousDataInDbIsRemoved(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	typesToSave := GetEntriesTypes()
	err := dataProvider.SaveEntriesTypesToDb(typesToSave)
	if err != nil {
		log.Fatal(err)
	}
	err = dataProvider.SaveEntriesTypesToDb([]EntryType{})
	if err != nil {
		log.Fatal(err)
	}
	loadedTypes, err := dataProvider.LoadEntriesTypesFromDb()
	assert.Empty(t, loadedTypes)
	DeleteTestDb()
}
