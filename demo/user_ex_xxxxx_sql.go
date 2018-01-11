// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.

package demo

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/rickb777/sqlgen2"
	"github.com/rickb777/sqlgen2/schema"
	"log"
)

// XUserTable holds a given table name with the database reference, providing access methods below.
// The Prefix field is often blank but can be used to hold a table name prefix (e.g. ending in '_'). Or it can
// specify the name of the schema, in which case it should have a trailing '.'.
type XUserTable struct {
	name    sqlgen2.TableName
	db      sqlgen2.Execer
	ctx     context.Context
	dialect schema.Dialect
	logger  *log.Logger
}

// Type conformance check
var _ sqlgen2.Table = &XUserTable{}

// NewXUserTable returns a new table instance.
// If a blank table name is supplied, the default name "users" will be used instead.
// The request context is initialised with the background.
func NewXUserTable(name sqlgen2.TableName, d sqlgen2.Execer, dialect schema.Dialect) XUserTable {
	if name.Name == "" {
		name.Name = "users"
	}
	return XUserTable{name, d, context.Background(), dialect, nil}
}

// CopyTableAsXUserTable copies a table instance, retaining the name etc but
// providing methods appropriate for 'User'.
func CopyTableAsXUserTable(origin sqlgen2.Table) XUserTable {
	return XUserTable{
		name:    origin.Name(),
		db:      origin.DB(),
		ctx:     origin.Ctx(),
		dialect: origin.Dialect(),
		logger:  origin.Logger(),
	}
}

// WithPrefix sets the table name prefix for subsequent queries.
func (tbl XUserTable) WithPrefix(pfx string) XUserTable {
	tbl.name.Prefix = pfx
	return tbl
}

// WithContext sets the context for subsequent queries.
func (tbl XUserTable) WithContext(ctx context.Context) XUserTable {
	tbl.ctx = ctx
	return tbl
}

// WithLogger sets the logger for subsequent queries.
func (tbl XUserTable) WithLogger(logger *log.Logger) XUserTable {
	tbl.logger = logger
	return tbl
}

// Ctx gets the current request context.
func (tbl XUserTable) Ctx() context.Context {
	return tbl.ctx
}

// Dialect gets the database dialect.
func (tbl XUserTable) Dialect() schema.Dialect {
	return tbl.dialect
}

// Logger gets the trace logger.
func (tbl XUserTable) Logger() *log.Logger {
	return tbl.logger
}

// SetLogger sets the logger for subsequent queries, returning the interface.
func (tbl XUserTable) SetLogger(logger *log.Logger) sqlgen2.Table {
	tbl.logger = logger
	return tbl
}

// Name gets the table name.
func (tbl XUserTable) Name() sqlgen2.TableName {
	return tbl.name
}

// DB gets the wrapped database handle, provided this is not within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl XUserTable) DB() *sql.DB {
	return tbl.db.(*sql.DB)
}

// Tx gets the wrapped transaction handle, provided this is within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl XUserTable) Tx() *sql.Tx {
	return tbl.db.(*sql.Tx)
}

// IsTx tests whether this is within a transaction.
func (tbl XUserTable) IsTx() bool {
	_, ok := tbl.db.(*sql.Tx)
	return ok
}

// Begin starts a transaction. The default isolation level is dependent on the driver.
func (tbl XUserTable) BeginTx(opts *sql.TxOptions) (XUserTable, error) {
	d := tbl.db.(*sql.DB)
	var err error
	tbl.db, err = d.BeginTx(tbl.ctx, opts)
	return tbl, err
}

func (tbl XUserTable) logQuery(query string, args ...interface{}) {
	sqlgen2.LogQuery(tbl.logger, query, args...)
}


//--------------------------------------------------------------------------------

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// It returns the number of rows affected (of the database drive supports this).
func (tbl XUserTable) Exec(query string, args ...interface{}) (int64, error) {
	tbl.logQuery(query, args...)
	res, err := tbl.db.ExecContext(tbl.ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//--------------------------------------------------------------------------------

// QueryOne is the low-level access function for one User.
// If the query selected many rows, only the first is returned; the rest are discarded.
// If not found, *User will be nil.
func (tbl XUserTable) QueryOne(query string, args ...interface{}) (*User, error) {
	list, err := tbl.doQuery(true, query, args...)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return list[0], nil
}

// Query is the low-level access function for Users.
func (tbl XUserTable) Query(query string, args ...interface{}) ([]*User, error) {
	return tbl.doQuery(false, query, args...)
}

func (tbl XUserTable) doQuery(firstOnly bool, query string, args ...interface{}) ([]*User, error) {
	tbl.logQuery(query, args...)
	rows, err := tbl.db.QueryContext(tbl.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanXUsers(rows, firstOnly)
}

// scanXUsers reads table records into a slice of values.
func scanXUsers(rows *sql.Rows, firstOnly bool) ([]*User, error) {
	var err error
	var vv []*User

	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 bool
	var v5 bool
	var v6 []byte
	var v7 int64
	var v8 string
	var v9 string

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
		)
		if err != nil {
			return vv, err
		}

		v := &User{}
		v.Uid = v0
		v.Login = v1
		v.EmailAddress = v2
		v.Avatar = v3
		v.Active = v4
		v.Admin = v5
		err = json.Unmarshal(v6, &v.Fave)
		if err != nil {
			return nil, err
		}
		v.LastUpdated = v7
		v.token = v8
		v.secret = v9

		var iv interface{} = v
		if hook, ok := iv.(sqlgen2.CanPostGet); ok {
			err = hook.PostGet()
			if err != nil {
				return vv, err
			}
		}

		vv = append(vv, v)

		if firstOnly {
			return vv, rows.Err()
		}
	}

	return vv, rows.Err()
}