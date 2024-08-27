package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Test some happy paths
func TestPGStoreCRUD(t *testing.T) {
	store, close := prepareEmptyTestDB(t)
	defer close()
	acc, err := NewAccount("firstName", "lastname", "secret_password")

	// list available accounts
	accounts, err := store.GetAccounts()
	assert.Nil(t, err)
	assert.Len(t, *accounts, 0, "This should be empty")

	// create an account
	assert.Nil(t, err)
	assert.Nil(t, store.CreateAccount(acc))
	assert.NotEqual(t, 0, acc.ID, "id should be set by now")

	// fetch the acount back by id
	getAccount, err := store.GetAccountByID(acc.ID)
	assert.Nil(t, err)
	assert.Equal(t, acc, getAccount)

	// fetch the acount back by id
	getAccount, err = store.GetAccountByNumber(int(acc.Number))
	assert.Nil(t, err)
	assert.Equal(t, acc, getAccount)

	// list available accounts
	accounts, err = store.GetAccounts()
	assert.Nil(t, err)
	assert.Len(t, *accounts, 1, "This should not be empty")

	// Update the account TODO

	// Delete the account
	err = store.DeleteAccount(acc.ID)
	assert.Nil(t, err)
	accounts, err = store.GetAccounts()
	assert.Nil(t, err)
	assert.Len(t, *accounts, 0, "This should be empty")

}

func prepareEmptyTestDB(t *testing.T) (storage Storage, close func() error) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	dbConfig := PostgresStoreConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		port:     os.Getenv("POSTGRES_PORT"),
		dbName:   os.Getenv("POSTGRES_DB"),
	}
	defaultStore, err := NewPostgressStore(&dbConfig)
	assert.Nil(t, err)
	testDB := dbConfig.dbName + "_test"
	_, err = defaultStore.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDB))
	assert.Nil(t, err)
	_, err = defaultStore.db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDB))
	assert.Nil(t, err)
	defaultStore.db.Close()
	testdbConfig := PostgresStoreConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		port:     os.Getenv("POSTGRES_PORT"),
		dbName:   testDB,
	}
	store, err := NewPostgressStore(&testdbConfig)
	assert.Nil(t, err)
	store.Init()
	assert.Nil(t, err)
	close = func() error {
		fmt.Println("Closing Test DB connection")
		if err := store.db.Close(); err != nil {
			fmt.Println("!!!!!!! FAILED to close DB !!!!!!!")
		}
		if err := store.db.Ping(); err != nil {
			fmt.Println("connection close... ping error")
		} else {
			fmt.Println("connection NOOOOOT close... ping success")
		}
		return err

	}

	return store, close

}
