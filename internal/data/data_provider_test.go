package data

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

const testDbPath = "../../test/testDb.db"

func deleteOldTestDb() {
	_, err := os.Stat(testDbPath)
	if os.IsExist(err) {
		err = os.Remove(testDbPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getTestEntriesToSave() []Entry {
	return []Entry{
		{
			Id:                         0,
			Status:                     "some status",
			Title:                      "some title",
			Completion:                 1,
			AmountOfElementsToComplete: 2,
			Score:                      3,
			Link:                       "some link",
			Description:                "some description",
			Comment:                    "some comment",
			MediaType:                  "some media type",
		},
		{
			Id:                         1,
			Status:                     "some status2",
			Title:                      "some title2",
			Completion:                 4,
			AmountOfElementsToComplete: 5,
			Score:                      6,
			Link:                       "some link2",
			Description:                "some description2",
			Comment:                    "some comment2",
			MediaType:                  "some media type2",
		},
	}
}

func TestDbOperations(t *testing.T) {
	deleteOldTestDb()
	entriesToSave := getTestEntriesToSave()
	dataProvider := newDataProvider(testDbPath)
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
