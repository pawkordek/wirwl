package data

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const entriesTypesTableName = "entries_types"
const entriesTableSuffix = "_entries"

type BoltProvider struct {
	dbPath string
	db     *bolt.DB
}

func NewBoltProvider(dbPath string) Provider {
	return &BoltProvider{dbPath: dbPath}
}

func (provider *BoltProvider) SaveEntries(entries map[EntryType][]Entry) error {
	entriesTypes := provider.getEntriesTypesFromEntries(entries)
	err := provider.saveEntriesTypesToDb(entriesTypes)
	if err != nil {
		return err
	}
	for entryType, entries := range entries {
		err = provider.saveEntriesToDb(entryType.Name+entriesTableSuffix, entries)
		if err != nil {
			return err
		}
	}
	return nil
}

func (provider *BoltProvider) getEntriesTypesFromEntries(entries map[EntryType][]Entry) []EntryType {
	var entriesTypes []EntryType
	for entryType, _ := range entries {
		entriesTypes = append(entriesTypes, entryType)
	}
	return entriesTypes
}

func (provider *BoltProvider) openDb() error {
	db, err := bolt.Open(provider.dbPath, 0600, &bolt.Options{Timeout: 10 * time.Second})
	provider.db = db
	if err != nil {
		return errors.Wrap(err, "An error occurred when opening the database")
	}
	return nil
}

func (provider *BoltProvider) closeDb() error {
	err := provider.db.Close()
	if err != nil {
		return errors.Wrap(err, "An error occurred when closing the database")
	}
	return nil
}

func (provider *BoltProvider) saveEntriesToDb(table string, entries []Entry) error {
	err := provider.openDb()
	if err != nil {
		return err
	}
	defer func() {
		err = provider.closeDb()
	}()
	err = provider.saveEntriesInTable(table, entries)
	if err != nil {
		return err
	}
	return err
}

func (provider *BoltProvider) saveEntriesInTable(table string, entries []Entry) error {
	err := provider.deleteTableIfExists(table)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		err = provider.createNewTable(table)
		if err != nil {
			return err
		}
	} else {
		for _, entry := range entries {
			err := provider.saveEntryToTable(table, entry)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (provider *BoltProvider) createNewTable(name string) error {
	return provider.db.Update(func(transaction *bolt.Tx) error {
		_, err := transaction.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return errors.Wrap(err, "An error occurred when creating a new table with name "+name)
		}
		return err
	})
}

func (provider *BoltProvider) saveEntryToTable(table string, entry Entry) error {
	entryAsJSON, err := json.Marshal(entry)
	if err != nil {
		return errors.Wrap(err, "An error occurred when marshaling entry during entry saving to table with name "+table+". Entry to save was: "+entry.String())
	}
	err = provider.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte(table))
		if err != nil {
			return errors.Wrap(err, "An error occurred when creating bucket for table with name "+table+" during entry saving. Entry to save was: "+entry.String())
		}
		return bucket.Put([]byte(strconv.Itoa(entry.Id)), entryAsJSON)
	})
	if err != nil {
		return errors.Wrap(err, "An error occurred when making update on the database during entry saving for table with name "+table+". Entry to save was:"+entry.String())
	}
	return nil
}

func (provider *BoltProvider) deleteTableIfExists(table string) error {
	return provider.db.Update(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket([]byte(table))
		if bucket != nil {
			err := transaction.DeleteBucket([]byte(table))
			if err != nil {
				return errors.Wrap(err, "An error occurred when deleting an existing table with name "+table)
			}
		}
		return nil
	})
}

func (provider *BoltProvider) LoadEntries() (map[EntryType][]Entry, error) {
	entriesTypes, err := provider.loadEntriesTypesFromDb()
	if err != nil {
		return nil, err
	}
	allEntries := make(map[EntryType][]Entry)
	for _, entryType := range entriesTypes {
		entriesForCurrentType, err := provider.loadEntriesFromDb(entryType.Name + entriesTableSuffix)
		if err != nil {
			return nil, err
		}
		allEntries[entryType] = entriesForCurrentType
	}
	return allEntries, nil
}

func (provider *BoltProvider) loadEntriesFromDb(table string) ([]Entry, error) {
	err := provider.openDb()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = provider.closeDb()
	}()
	var entries []Entry
	err = provider.db.View(func(transaction *bolt.Tx) error {
		entries, err = getEntriesDataFromTable(transaction, table)
		return err
	})
	if err != nil {
		return nil, err
	}
	return entries, err
}

func getEntriesDataFromTable(transaction *bolt.Tx, table string) ([]Entry, error) {
	var entries []Entry
	bucket := transaction.Bucket([]byte(table))
	if bucket == nil {
		return nil, errors.New("An error occurred when loading entries from table with name " + table + ". No such table")
	}
	err := bucket.ForEach(func(key, value []byte) error {
		var entry Entry
		err := json.Unmarshal(value, &entry)
		if err != nil {
			return errors.Wrap(err, "An error occurred when unmarshalling an entry "+entry.String())
		}
		entries = append(entries, entry)
		return nil
	})
	return entries, err
}

func (provider *BoltProvider) saveEntriesTypesToDb(entriesTypes []EntryType) error {
	err := provider.openDb()
	if err != nil {
		return err
	}
	err = provider.deleteTableIfExists(entriesTypesTableName)
	for _, entryType := range entriesTypes {
		err := provider.saveEntryTypeToTable(entryType)
		if err != nil {
			return err
		}
	}
	return provider.db.Close()
}

func (provider *BoltProvider) saveEntryTypeToTable(entryType EntryType) error {
	typeAsJSON, err := json.Marshal(entryType)
	if err != nil {
		return errors.Wrap(err, "An error occurred when marshaling entry type during entry type saving")
	}
	err = provider.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte(entriesTypesTableName))
		if err != nil {
			return errors.Wrap(err, "An error occurred when creating bucket during entry type saving")
		}
		return bucket.Put([]byte(entryType.Name), typeAsJSON)
	})
	if err != nil {
		return errors.Wrap(err, "An error occurred when making update on the database during entry type saving")
	}
	return nil
}

func (provider *BoltProvider) loadEntriesTypesFromDb() ([]EntryType, error) {
	err := provider.openDb()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = provider.closeDb()
	}()
	var types []EntryType
	err = provider.db.View(func(transaction *bolt.Tx) error {
		types, err = getEntriesTypesFromTable(transaction)
		return err
	})
	if err != nil {
		return nil, err
	}
	return types, err
}

func getEntriesTypesFromTable(transaction *bolt.Tx) ([]EntryType, error) {
	var types []EntryType
	bucket := transaction.Bucket([]byte(entriesTypesTableName))
	if bucket == nil {
		//If there is no entries types table it simply means that the db is new and types are yet to be created
		return types, nil
	}
	err := bucket.ForEach(func(key, value []byte) error {
		var entryType EntryType
		err := json.Unmarshal(value, &entryType)
		if err != nil {
			return errors.Wrap(err, "An error occurred when unmarshalling an entry type "+entryType.String())
		}
		types = append(types, entryType)
		return nil
	})
	return types, err
}
