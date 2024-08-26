package appServices

import (
	"encoding/json"
	"errors"
	"i9rfs/client/helpers"
	"io/fs"
	"os"
)

var storageFile string

type AppDataStore struct {
	storage map[string]any
}

func (adst AppDataStore) GetItem(key string, dest any) {
	helpers.ParseTo(adst.storage[key], dest)
}

func (adst *AppDataStore) SetItem(key string, value any) {
	adst.storage[key] = value
}

func (adst *AppDataStore) RemoveItem(key string) {
	delete(adst.storage, key)
}

func (adst AppDataStore) Save() {
	apstData, err := json.Marshal(adst.storage)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(storageFile, apstData, fs.ModePerm); err != nil {
		panic(err)
	}
}

func (adst *AppDataStore) Revive(inStorageFile string) error {
	apstData, err := os.ReadFile(inStorageFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(err)
	}

	storageFile = inStorageFile

	var store = make(map[string]any)

	if len(apstData) > 1 {
		json.Unmarshal(apstData, &store)
	}

	adst.storage = store

	return nil
}
