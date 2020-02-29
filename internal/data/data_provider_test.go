package data

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDbOperationsOnEntries(t *testing.T) {
	entriesToSave := GetTestEntriesToSave()
	dataProvider := NewDataProvider(TestDbPath)
	err := dataProvider.SaveEntriesToDb("test_table", entriesToSave)
	if err != nil {
		log.Fatal(err)
	}
	entries, err := dataProvider.LoadEntriesFromDb("test_table")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, entries[0], entriesToSave[0])
	assert.Equal(t, entries[1], entriesToSave[1])
	DeleteTestDb()
}

func TestDbOperationsOnEntriesTypes(t *testing.T) {
	entriesTypes := GetEntriesTypes()
	dataProvider := NewDataProvider(TestDbPath)
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
	DeleteTestDb()
}

func TestThatTryingToLoadEntriesFromEmptyDbReturnsEmptySlice(t *testing.T) {
	dataProvider := NewDataProvider(TestDbPath)
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
