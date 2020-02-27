package data

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
)

type Entry struct {
	Id                         int
	Status                     string
	Title                      string
	Completion                 int
	AmountOfElementsToComplete int
	Score                      int
	Link                       string
	Description                string
	Comment                    string
	MediaType                  string
}

type dataProvider struct {
	dbPath string
	db     *bolt.DB
}

func newDataProvider(dbPath string) *dataProvider {
	return &dataProvider{dbPath: dbPath}
}

func (provider *dataProvider) openDb() error {
	db, err := bolt.Open(provider.dbPath, 0600, nil)
	provider.db = db
	return err
}

func (provider *dataProvider) closeDb() error {
	return provider.db.Close()
}

func (provider *dataProvider) SaveEntriesToDb(table string, entries []Entry) error {
	err := provider.openDb()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err := provider.saveEntryToTable(table, entry)
		if err != nil {
			return err
		}
	}
	err = provider.closeDb()
	if err != nil {
		return err
	}
	return nil
}

func (provider *dataProvider) saveEntryToTable(table string, entry Entry) error {
	entryAsJSON, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	err = provider.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(strconv.Itoa(entry.Id)), entryAsJSON)
	})
	if err != nil {
		return err
	}
	return nil
}

func (provider *dataProvider) LoadEntriesFromDb(table string) ([]Entry, error) {
	err := provider.openDb()
	if err != nil {
		return nil, err
	}
	var entries [] Entry
	err = provider.db.View(func(transaction *bolt.Tx) error {
		entries, err = getDbDataLoadedIntoSlice(transaction, table)
		return err
	})
	if err != nil {
		return nil, err
	}
	err = provider.closeDb()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func getDbDataLoadedIntoSlice(transaction *bolt.Tx, table string) ([]Entry, error) {
	var entries []Entry
	bucket := transaction.Bucket([]byte(table))
	if bucket == nil {
		return nil, errors.New("No table with name=" + table)
	}
	err := bucket.ForEach(func(key, value []byte) error {
		var entry Entry
		err := json.Unmarshal(value, &entry)
		if err != nil {
			return err
		}
		entries = append(entries, entry)
		return nil
	})
	return entries, err
}
