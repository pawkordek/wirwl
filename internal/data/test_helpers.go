package data

import (
	"log"
	"os"
)

const ExampleDbPath = "../../test/exampleDb.db"
const TestDbPath = "../../test/testDb.db"
const EmptyDbPath = "../../test/emptyDb.db"

func GetEntriesTypes() []string {
	return []string{
		"comic",
		"video",
	}
}

func GetExampleVideoEntries() []Entry {
	return []Entry{
		{
			Id:                         0,
			Status:                     "some status",
			Title:                      "some video1",
			Completion:                 1,
			AmountOfElementsToComplete: 2,
			Score:                      3,
			Link:                       "some link",
			Description:                "some description",
			Comment:                    "some comment",
			MediaType:                  "video",
		},
		{
			Id:                         1,
			Status:                     "some status2",
			Title:                      "some video2",
			Completion:                 4,
			AmountOfElementsToComplete: 5,
			Score:                      6,
			Link:                       "some link2",
			Description:                "some description2",
			Comment:                    "some comment2",
			MediaType:                  "video",
		},
	}
}

func GetExampleComicEntries() []Entry {
	return []Entry{
		{
			Id:                         0,
			Status:                     "some status",
			Title:                      "some comic1",
			Completion:                 1,
			AmountOfElementsToComplete: 2,
			Score:                      3,
			Link:                       "some link",
			Description:                "some description",
			Comment:                    "some comment",
			MediaType:                  "comic",
		},
		{
			Id:                         1,
			Status:                     "some status2",
			Title:                      "some comic2",
			Completion:                 4,
			AmountOfElementsToComplete: 5,
			Score:                      6,
			Link:                       "some link2",
			Description:                "some description2",
			Comment:                    "some comment2",
			MediaType:                  "comic",
		},
	}
}

func DeleteTestDb() {
	DeleteFile(TestDbPath)
}

func DeleteFile(path string) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetTestEntriesToSave() []Entry {
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
