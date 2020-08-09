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

func TestThatTryingToLoadEntriesFromEmptyDbReturnsEmptySlice(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	entries, err := dataProvider.LoadEntries()
	assert.Equal(t, 0, len(entries))
	assert.Nil(t, err)
	DeleteTestDb()
}

func TestThatSavingEmptyEntriesDoesNotThrowError(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	err := dataProvider.SaveEntries(map[EntryType][]Entry{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = dataProvider.LoadEntries()
	assert.Nil(t, err)
	DeleteTestDb()
}

func TestThatWhenSavingEntriesPreviousDataInDbIsRemoved(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	dataProvider := NewBoltProvider(testDbPath)
	entriesToSave := GetTestEntries()
	err := dataProvider.SaveEntries(entriesToSave)
	if err != nil {
		log.Fatal(err)
	}
	err = dataProvider.SaveEntries(map[EntryType][]Entry{})
	if err != nil {
		log.Fatal(err)
	}
	loadedEntries, err := dataProvider.LoadEntries()
	assert.Empty(t, loadedEntries)
	DeleteTestDb()
}
