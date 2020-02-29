package data

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDbOperationsOnEntries(t *testing.T) {
	DeleteOldTestDb()
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
}

func TestDbOperationsOnEntriesTypes(t *testing.T) {
	DeleteOldTestDb()
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
}
