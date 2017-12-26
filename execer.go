package sqlgen2

import (
	"database/sql"
	"context"
	"github.com/rickb777/sqlgen2/where"
)

// Execer describes the methods of the core database API. See database/sql/DB and database/sql/Tx.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Type conformance assertions
var _ Execer = &sql.DB{}
var _ Execer = &sql.Tx{}

type CanPreInsert interface {
	PreInsert(Execer) error
}

type CanPreUpdate interface {
	PreUpdate(Execer) error
}

//-------------------------------------------------------------------------------------------------

// Table provides the generic features of each generated table handler.
type Table interface {
	// FullName gets the concatenated prefix and table name.
	FullName() string

	// DB gets the wrapped database handle, provided this is not within a transaction.
	// Panics if it is in the wrong state - use IsTx() if necessary.
	DB() *sql.DB

	// Tx gets the wrapped transaction handle, provided this is within a transaction.
	// Panics if it is in the wrong state - use IsTx() if necessary.
	Tx() *sql.Tx

	// IsTx tests whether this is within a transaction.
	IsTx() bool
}

type TableCreator interface {
	Table

	// CreateTable creates the database table using the current dialect.
	CreateTable(ifNotExist bool) (int64, error)
}

type TableWithIndexes interface {
	TableCreator

	// CreateIndexes creates the indexes for the database table using the current dialect.
	CreateIndexes(ifNotExist bool) (err error)

	// CreateTableWithIndexes creates the database table and its indexes using the current dialect.
	CreateTableWithIndexes(ifNotExist bool) (err error)
}

type TableWithCrud interface {
	Table

	// Exec executes a query.
	// It returns the number of rows affected (if the DB supports that).
	Exec(query string, args ...interface{}) (int64, error)

	// CountSA counts records that match a 'where' predicate.
	CountSA(where string, args ...interface{}) (count int64, err error)

	// Count counts records that match a 'where' predicate.
	Count(where where.Expression) (count int64, err error)

	// UpdateFields writes new values to the specified columns for rows that match the 'where' predicate.
	// It returns the number of rows affected (if the DB supports that).
	UpdateFields(where where.Expression, fields ...sql.NamedArg) (int64, error)

	// Delete deletes rows that match the 'where' predicate.
	// It returns the number of rows affected (if the DB supports that).
	Delete(where where.Expression) (int64, error)

	// These methods are provided but have a specific record type so do not conform to a general interface.
	//QueryOne(query string, args ...interface{}) (*User, error)
	//Query(query string, args ...interface{}) ([]*User, error)
	//SelectOneSA(where, orderBy string, args ...interface{}) (*User, error)
	//SelectOne(where where.Expression, orderBy string) (*User, error)
	//SelectSA(where, orderBy string, args ...interface{}) ([]*User, error)
	//Select(where where.Expression, orderBy string) ([]*User, error)
	//Insert(vv ...*User) error
	//Update(vv ...*User) (int64, error)
}