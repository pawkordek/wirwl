package data

import (
	"errors"
)

const TestDbPath = "../../testdata/testDb.db"

func GetEntriesTypes() []EntryType {
	return []EntryType{
		{
			Name:       "comics",
			ImageQuery: "comic cover",
		},
		{
			Name:       "music",
			ImageQuery: "album cover",
		},
		{
			Name:       "videos",
			ImageQuery: "video cover",
		},
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
			Tags:                       "some tags",
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
			Tags:                       "some tags",
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
			Tags:                       "some tags",
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
			Tags:                       "some tags",
		},
	}
}

func GetExampleMusicEntries() []Entry {
	return []Entry{
		{
			Id:                         0,
			Status:                     "some status",
			Title:                      "some music1",
			Completion:                 1,
			AmountOfElementsToComplete: 2,
			Score:                      3,
			Link:                       "some link",
			Description:                "some description",
			Comment:                    "some comment",
			MediaType:                  "music",
			Tags:                       "some tags",
		},
		{
			Id:                         1,
			Status:                     "some status2",
			Title:                      "some music2",
			Completion:                 4,
			AmountOfElementsToComplete: 5,
			Score:                      6,
			Link:                       "some link2",
			Description:                "some description2",
			Comment:                    "some comment2",
			MediaType:                  "music",
			Tags:                       "some tags",
		},
	}
}

func DeleteTestDb() {
	DeleteDirWithContents(TestDbPath)
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
			Tags:                       "some tags",
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
			Tags:                       "some tags",
		},
	}
}

var AlwaysFailingProviderError = errors.New("Error from always failing provider")

type AlwaysFailingProvider struct {
}

func NewAlwaysFailingProvider() Provider {
	return &AlwaysFailingProvider{}
}

func (provider *AlwaysFailingProvider) SaveEntriesToDb(table string, entries []Entry) error {
	return AlwaysFailingProviderError
}

func (provider *AlwaysFailingProvider) LoadEntriesFromDb(table string) ([]Entry, error) {
	return nil, AlwaysFailingProviderError
}

func (provider *AlwaysFailingProvider) SaveEntriesTypesToDb(entriesTypes []EntryType) error {
	return AlwaysFailingProviderError
}

func (provider *AlwaysFailingProvider) LoadEntriesTypesFromDb() ([]EntryType, error) {
	return nil, AlwaysFailingProviderError
}
