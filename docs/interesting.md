# Some interesting findings

## Go missing some SQL parameter escaping features, like identifier paramaters.

GitHub issue: [https://github.com/golang/go/issues/18478](https://github.com/golang/go/issues/18478)

If we want to modify the following query to be able to create `account_test` table, cannot bind this value.

```go
func (s *PostgresStore) addEncryptedPasswordField() error {
	query := "ALTER TABLE account ADD COLUMN IF NOT EXISTS encrypted_password varchar(100)"
	_, err := s.db.Exec(query)
	return err
}
```

```go
func (s *PostgresStore) addEncryptedPasswordField() error {
	tableName := "account_test"
    query := "ALTER TABLE $1 ADD COLUMN IF NOT EXISTS encrypted_password varchar(100)"
	_, err := s.db.Exec(query, tableName)
	return err
}
```

We could use `fmt.Sprintf`, but it is not a "prepared statment". Be careful.

```go
func (s *PostgresStore) addEncryptedPasswordField() error {
	tableName := "account_test"
    query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS encrypted_password varchar(100)",tableName)
	_, err := s.db.Exec(query)
	return err
}
```

## Integration test

Mocking db connection is great, but some times, we want to store and fetch the data. :)
With integration test we can find issus like:

```
pg_storage_test.go:22:
                Error Trace:    /gobank/pg_storage_test.go:22
                Error:          Not equal:
                                expected: &main.Account{ID:1, FirstName:"firstName", LastName:"lastname", EncryptedPassword:"hash", Number:4140261, Balance:0, CreatedAt:time.Date(2024, time.August, 26, 16, 31, 15, 784850000, time.UTC)}
                                actual  : &main.Account{ID:1, FirstName:"firstName", LastName:"lastname", EncryptedPassword:"hash", Number:4140261, Balance:0, CreatedAt:time.Date(2024, time.August, 26, 16, 31, 15, 784850000, time.Location(""))}
```

## E2E test

I know I could just use the httptest server and pass it some handler functions to make some meaningful E2E tests.

But this blackbox like approach gived me the opportinuty to tryout something different. And it also revealed and interesting bug:

```bash
Error:          Expected nil, but got: &pq.Error{Severity:"ERROR", Code:"55006", Message:"database \"postgres_test\" is being accessed by other users", Detail:"There are 3 other sessions using the database.", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"dbcommands.c", Line:"1699", Routine:"dropdb"}
```

Explicitly closing the db connection between testcases to create a clean test database failed, because in the tutorial rows were not closed after processing them. The `go vet` command did not revealed any issues, but the test stil failed.

GitHub issue: [https://github.com/golang/go/issues/34544](https://github.com/golang/go/issues/34544)

The issue is closed, but I think it still need some adjustments.
