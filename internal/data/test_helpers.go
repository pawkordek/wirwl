package data

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"wirwl/internal/log"
)

const TestDbPath = "../../testdata/testDb.db"

var comicsEntryType = EntryType{
	Name:       "comics",
	ImageQuery: "comic cover",
}

var musicEntryType = EntryType{
	Name:       "music",
	ImageQuery: "album cover",
}

var videoEntryType = EntryType{
	Name:       "videos",
	ImageQuery: "video cover",
}

func getTempDbPath() (string, func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "test.db"), func() { _ = os.RemoveAll(dir) }
}

func GetTestEntries() map[EntryType][]Entry {
	entries := make(map[EntryType][]Entry)
	entries[comicsEntryType] = GetExampleComicEntries()
	entries[musicEntryType] = GetExampleMusicEntries()
	entries[videoEntryType] = GetExampleVideoEntries()
	return entries
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

func (provider *AlwaysFailingProvider) SaveEntries(map[EntryType][]Entry) error {
	return AlwaysFailingProviderError
}

func (provider *AlwaysFailingProvider) LoadEntries() (map[EntryType][]Entry, error) {
	return nil, AlwaysFailingProviderError
}

/*It is supposed to provide default functions that return empty values but every single one can be overwritten
so that desired functionality when testing can be achieved
*/
type AbstractProvider struct {
	SaveEntriesFunc            func(map[EntryType][]Entry) error
	LoadEntriesFunc            func() (map[EntryType][]Entry, error)
}

func NewAbstractProvider() *AbstractProvider {
	return &AbstractProvider{
		SaveEntriesFunc: func(entries map[EntryType][]Entry) error {
			return nil
		},
		LoadEntriesFunc: func() (map[EntryType][]Entry, error) {
			return make(map[EntryType][]Entry), nil
		},
	}
}

func (provider *AbstractProvider) SaveEntries(entries map[EntryType][]Entry) error {
	return provider.SaveEntriesFunc(entries)
}

func (provider *AbstractProvider) LoadEntries() (map[EntryType][]Entry, error) {
	return provider.LoadEntriesFunc()
}
