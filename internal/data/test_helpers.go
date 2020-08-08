package data

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"wirwl/internal/log"
)

const TestDbPath = "../../testdata/testDb.db"

func getTempDbPath() (string, func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "test.db"), func() { _ = os.RemoveAll(dir) }
}

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
	err := DeleteDirWithContents(TestDbPath)
	if err != nil {
		log.Fatal(err)
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

/*It is supposed to provide default functions that return empty values but every single one can be overwritten
so that desired functionality when testing can be achieved
*/
type AbstractProvider struct {
	SaveEntriesToDbFunc        func(string, []Entry) error
	LoadEntriesFromDbFunc      func(string) ([]Entry, error)
	SaveEntriesTypesToDbFunc   func([]EntryType) error
	LoadEntriesTypesFromDbFunc func() ([]EntryType, error)
}

func NewAbstractProvider() *AbstractProvider {
	return &AbstractProvider{
		SaveEntriesToDbFunc: func(table string, entries []Entry) error {
			return nil
		},
		LoadEntriesFromDbFunc: func(s string) (entries []Entry, err error) {
			return []Entry{}, nil
		},
		SaveEntriesTypesToDbFunc: func(types []EntryType) error {
			return nil
		},
		LoadEntriesTypesFromDbFunc: func() (types []EntryType, err error) {
			return []EntryType{}, nil
		},
	}
}

func (provider *AbstractProvider) SaveEntriesToDb(table string, entries []Entry) error {
	return provider.SaveEntriesToDbFunc(table, entries)
}

func (provider *AbstractProvider) LoadEntriesFromDb(table string) ([]Entry, error) {
	return provider.LoadEntriesFromDbFunc(table)
}

func (provider *AbstractProvider) SaveEntriesTypesToDb(entriesTypes []EntryType) error {
	return provider.SaveEntriesTypesToDbFunc(entriesTypes)
}

func (provider *AbstractProvider) LoadEntriesTypesFromDb() ([]EntryType, error) {
	return provider.LoadEntriesTypesFromDbFunc()
}
