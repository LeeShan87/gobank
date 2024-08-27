package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	GetAccounts() (*[]Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

type PostgresStoreConfig struct {
	user     string
	password string
	dbName   string
	port     string
}

func NewPostgressStore(config *PostgresStoreConfig) (*PostgresStore, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		config.user,
		config.password,
		config.port,
		config.dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account(
	first_name, 
	last_name, 
	encrypted_password,
	number,
	balance, 
	created_at) 
	values ($1,$2,$3,$4,$5,$6) RETURNING id`
	if err := s.db.QueryRow(
		query,
		account.FirstName,
		account.LastName,
		account.EncryptedPassword,
		account.Number,
		account.Balance,
		account.CreatedAt).Scan(&account.ID); err != nil {
		return err
	}
	fmt.Println("Create account: ", account)
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := "SELECT * FROM account WHERE id = $1"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("account %d not found", id)
	}
	account, err := scanIntoAccount(rows)
	return &account, err

}

func (s *PostgresStore) GetAccountByNumber(id int) (*Account, error) {
	query := "SELECT * FROM account WHERE number = $1"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, fmt.Errorf("account with number [%d] not found", id)
	}
	account, err := scanIntoAccount(rows)
	return &account, err

}

func (s *PostgresStore) GetAccounts() (*[]Account, error) {
	query := "SELECT * FROM account"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return &accounts, err
		}
		accounts = append(accounts, account)
	}

	return &accounts, nil
}

func (s *PostgresStore) UpdateAccount(acount *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := "DELETE FROM account WHERE id = $1"
	_, err := s.db.Exec(query, id)

	return err
}

func scanIntoAccount(rows *sql.Rows) (Account, error) {
	var account Account
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.EncryptedPassword,
	)
	if err != nil {
		return account, err
	}

	// Convert CreatedAt to UTC
	account.CreatedAt = account.CreatedAt.UTC()
	return account, err
}
