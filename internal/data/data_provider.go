package data

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
)

const entriesTypesTableName = "entries_types"

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

type Provider interface {
	SaveEntriesToDb(table string, entries []Entry) error
	LoadEntriesFromDb(table string) ([]Entry, error)
	SaveEntriesTypesToDb(entriesTypes []string) error
	LoadEntriesTypesFromDb() ([]string, error)
}

type BoltProvider struct {
	dbPath string
	db     *bolt.DB
}

func NewBoltProvider(dbPath string) Provider {
	return &BoltProvider{dbPath: dbPath}
}

func (provider *BoltProvider) openDb() error {
	db, err := bolt.Open(provider.dbPath, 0600, nil)
	provider.db = db
	return err
}

func (provider *BoltProvider) closeDb() error {
	return provider.db.Close()
}

func (provider *BoltProvider) SaveEntriesToDb(table string, entries []Entry) error {
	err := provider.openDb()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		err = provider.createNewTable(table)
	} else {
		for _, entry := range entries {
			err := provider.saveEntryToTable(table, entry)
			if err != nil {
				return err
			}
		}
	}
	err = provider.closeDb()
	if err != nil {
		return err
	}
	return nil
}

func (provider *BoltProvider) createNewTable(name string) error {
	return provider.db.Update(func(transaction *bolt.Tx) error {
		_, err := transaction.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func (provider *BoltProvider) saveEntryToTable(table string, entry Entry) error {
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

func (provider *BoltProvider) LoadEntriesFromDb(table string) ([]Entry, error) {
	err := provider.openDb()
	if err != nil {
		return nil, err
	}
	var entries [] Entry
	err = provider.db.View(func(transaction *bolt.Tx) error {
		entries, err = getEntriesDataFromTable(transaction, table)
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

func getEntriesDataFromTable(transaction *bolt.Tx, table string) ([]Entry, error) {
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

func (provider *BoltProvider) SaveEntriesTypesToDb(entriesTypes []string) error {
	err := provider.openDb()
	if err != nil {
		return err
	}
	for _, entryType := range entriesTypes {
		err := provider.saveEntryTypeToTable(entryType)
		if err != nil {
			return err
		}
	}
	return provider.db.Close()
}

func (provider *BoltProvider) saveEntryTypeToTable(entryType string) error {
	return provider.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte(entriesTypesTableName))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(entryType), []byte(entryType))
	})
}

func (provider *BoltProvider) LoadEntriesTypesFromDb() ([]string, error) {
	err := provider.openDb()
	if err != nil {
		return nil, err
	}
	var types []string
	err = provider.db.View(func(transaction *bolt.Tx) error {
		types, err = getEntriesTypesFromTable(transaction)
		return err
	})
	if err != nil {
		return nil, err
	}
	err = provider.closeDb()
	if err != nil {
		return nil, err
	}
	return types, nil
}

func getEntriesTypesFromTable(transaction *bolt.Tx) ([]string, error) {
	var types []string
	bucket := transaction.Bucket([]byte(entriesTypesTableName))
	if bucket == nil {
		return types, nil
	}
	err := bucket.ForEach(func(key, value []byte) error {
		types = append(types, string(key))
		return nil
	})
	return types, err
}
