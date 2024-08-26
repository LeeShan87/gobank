package main

import "fmt"

func (s *PostgresStore) Init() error {
	if err := s.createAccountTable(); err != nil {
		return err
	}
	if err := s.addEncryptedPasswordField(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) dropTableForSeed() error {
	fmt.Println("Drop account table")
	query := "DROP TABLE IF EXISTS account"
	_, err := s.db.Exec(query)
	return err
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

func (s *PostgresStore) addEncryptedPasswordField() error {
	query := "ALTER TABLE account ADD COLUMN IF NOT EXISTS encrypted_password varchar(100)"
	_, err := s.db.Exec(query)
	return err
}
