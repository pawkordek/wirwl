package data

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDbOperationsOnEntries(t *testing.T) {
	entriesToSave := GetTestEntriesToSave()
	dataProvider := NewBoltProvider(TestDbPath)
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

func TestThatTryingToLoadDataIntoNonExistingTableReturnsError(t *testing.T) {
	dataProvider := NewBoltProvider(TestDbPath)
	types, err := dataProvider.LoadEntriesFromDb("non existing table")
	assert.Nil(t, types)
	assert.Equal(t, err, errors.New("No table with name=non existing table"))
	DeleteTestDb()
}

func TestDbOperationsOnEntriesTypes(t *testing.T) {
	entriesTypes := GetEntriesTypes()
	dataProvider := NewBoltProvider(TestDbPath)
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
	dataProvider := NewBoltProvider(TestDbPath)
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
	dataProvider := NewBoltProvider(TestDbPath)
	err := dataProvider.SaveEntriesToDb("new table", []Entry{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = dataProvider.LoadEntriesFromDb("new table")
	assert.Nil(t, err)
	DeleteTestDb()
}
