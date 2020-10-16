package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wirwl/internal/log"
)

func TestThatEntriesContainerHasDataAfterLoading(t *testing.T) {
	container := NewEntriesContainer(NewSampleTestDataProvider(""))
	err := container.LoadData()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, GetExampleComicEntries(), container.entries[comicsEntryType])
	assert.Equal(t, GetExampleVideoEntries(), container.entries[videoEntryType])
	assert.Equal(t, GetExampleMusicEntries(), container.entries[musicEntryType])
}

func TestThatEntriesContainerCreatesOutputDataFile(t *testing.T) {
	testDbPath, cleanup := getTempDbPath()
	defer cleanup()
	container := NewEntriesContainer(NewSampleTestDataProvider(testDbPath))
	err := container.SaveData()
	if err != nil {
		log.Fatal(err)
	}
	assert.FileExists(t, testDbPath)
}

func TestThatEntriesContainerReturnsAnErrorOnDataLoadFailure(t *testing.T) {
	container := NewEntriesContainer(NewAlwaysFailingProvider())
	err := container.LoadData()
	assert.NotNil(t, err)
}

func TestThatEntriesContainerReturnsAnErrorOnDataSaveFailuer(t *testing.T) {
	container := NewEntriesContainer(NewAlwaysFailingProvider())
	err := container.SaveData()
	assert.NotNil(t, err)
}

func TestThatAddingNewEntryTypeWorks(t *testing.T) {
	container := NewEntriesContainer(NewSampleTestDataProvider(""))
	typeToAdd := EntryType{
		Name:                  "added entry",
		CompletionElementName: "test element",
		ImageQuery:            "entry query",
	}
	_ = container.AddEntryType(typeToAdd)
	assert.NotNil(t, container.entries[typeToAdd])
}

func TestThatErrorIsReturnedWhenTryingToAddEntryTypeWithTheSameName(t *testing.T) {
	container := NewEntriesContainer(NewSampleTestDataProvider(""))
	typeToAdd := EntryType{
		Name:                  "added entry",
		CompletionElementName: "test element",
		ImageQuery:            "entry query",
	}
	typeToAdd2 := EntryType{
		Name:                  "added entry",
		CompletionElementName: "test element2",
		ImageQuery:            "entry query2",
	}
	_ = container.AddEntryType(typeToAdd)
	err := container.AddEntryType(typeToAdd2)
	assert.Contains(t, err.Error(), "Entry type with name added entry already exists")
}

func TestThatErrorIsReturnedWhenTryingToAddEntryTypeWithEmptyName(t *testing.T) {
	container := NewEntriesContainer(NewSampleTestDataProvider(""))
	typeToAdd := EntryType{
		Name:                  "",
		CompletionElementName: "test element",
		ImageQuery:            "entry query",
	}
	err := container.AddEntryType(typeToAdd)
	assert.Contains(t, err.Error(), "Cannot add entry type with an empty name")
}
