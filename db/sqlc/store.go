package db

import (
	"database/sql"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
}

// Store provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore Creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
