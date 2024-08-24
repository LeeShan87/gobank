package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountByID(int) error
	GetAccounts() (*[]Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "postgres://postgres:gobank@localhost:5400/postgres?sslmode=disable"
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

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}
func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account(
	first_name, 
	last_name, 
	number, 
	balance, 
	created_at) 
	values ($1,$2,$3,$4,$5) RETURNING id`
	if err := s.db.QueryRow(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt).Scan(&account.ID); err != nil {
		return err
	}
	fmt.Println("Create account: ", account)
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) error {
	return nil
}

func (s *PostgresStore) GetAccounts() (*[]Account, error) {
	query := "SELECT * FROM account"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for rows.Next() {
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			return &accounts, err
		}
		accounts = append(accounts, account)
	}

	return &accounts, nil
}
func (s *PostgresStore) UpdateAccount(acount *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}
