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
