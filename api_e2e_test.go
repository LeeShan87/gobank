package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// I know... I could do this just by testing the handlers
// But this is more fun :)
func TestAPI_E2E(t *testing.T) {
	db, close := prepareEmptyTestDB(t)
	defer close()
	server := NewApiServer("localhost:4444", db)

	go func() {
		server.Run()
	}()
	defer server.Shutdown()
	url := server.getListenAddress()
	assert.Equal(t, "http://localhost:4444", url)
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// accounts should be empty
	req, err := http.NewRequest(http.MethodGet, url+"/account", nil)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	var accounts []Account
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	assert.Nil(t, err)
	assert.Len(t, accounts, 0)

	// Create a new account
	payload := []byte(`
	{
    "firstName": "fname",
    "lastName": "lname",
    "password": "secret"
}
	`)
	req, err = http.NewRequest(http.MethodPost, url+"/account", bytes.NewReader(payload))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var acc Account
	err = json.NewDecoder(resp.Body).Decode(&acc)
	assert.Nil(t, err)
	assert.Equal(t, "fname", acc.FirstName)
	assert.Equal(t, "lname", acc.LastName)
	//assert.True(t, acc.validatePassword("secret"))

	// We should how have an account
	req, err = http.NewRequest(http.MethodGet, url+"/account", nil)
	assert.Nil(t, err)
	resp, err = client.Do(req)
	assert.Nil(t, err)
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	assert.Nil(t, err)
	assert.Len(t, accounts, 1)

	// Login with the account
	payload = []byte(fmt.Sprintf(`{
    "number" : %d,
    "password" : "%s"
}`, acc.Number, "secret"))
	req, err = http.NewRequest(http.MethodPost, url+"/login", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	resp, err = client.Do(req)
	assert.Nil(t, err)
	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	assert.Nil(t, err)
	assert.Equal(t, acc.Number, int64(loginResp.Number))

	// Get user by id (jwt protected :))
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/account/%d", url, acc.ID), nil)
	assert.Nil(t, err)
	// This converst header x-jwt-token to X-Jwt-Token
	// req.Header.Add("x-jwt-token", loginResp.Token)
	req.Header["x-jwt-token"] = []string{loginResp.Token}
	resp, err = client.Do(req)
	assert.Nil(t, err)
	var respAccount Account
	err = json.NewDecoder(resp.Body).Decode(&respAccount)
	assert.Nil(t, err)
	assert.Equal(t, acc.Number, respAccount.Number)

}
