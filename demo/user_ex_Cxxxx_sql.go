// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.

package demo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rickb777/sqlgen2"
	"github.com/rickb777/sqlgen2/schema"
	"log"
	"strings"
)

// CUserTable holds a given table name with the database reference, providing access methods below.
// The Prefix field is often blank but can be used to hold a table name prefix (e.g. ending in '_'). Or it can
// specify the name of the schema, in which case it should have a trailing '.'.
type CUserTable struct {
	name    sqlgen2.TableName
	db      sqlgen2.Execer
	ctx     context.Context
	dialect schema.Dialect
	logger  *log.Logger
	wrapper interface{}
}

// Type conformance check
var _ sqlgen2.Table = &CUserTable{}

// NewCUserTable returns a new table instance.
// If a blank table name is supplied, the default name "users" will be used instead.
// The request context is initialised with the background.
func NewCUserTable(name sqlgen2.TableName, d sqlgen2.Execer, dialect schema.Dialect) CUserTable {
	if name.Name == "" {
		name.Name = "users"
	}
	return CUserTable{name, d, context.Background(), dialect, nil, nil}
}

// CopyTableAsCUserTable copies a table instance, retaining the name etc but
// providing methods appropriate for 'User'.
func CopyTableAsCUserTable(origin sqlgen2.Table) CUserTable {
	return CUserTable{
		name:    origin.Name(),
		db:      origin.DB(),
		ctx:     origin.Ctx(),
		dialect: origin.Dialect(),
		logger:  origin.Logger(),
	}
}

// WithPrefix sets the table name prefix for subsequent queries.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) WithPrefix(pfx string) CUserTable {
	tbl.name.Prefix = pfx
	return tbl
}

// WithContext sets the context for subsequent queries.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) WithContext(ctx context.Context) CUserTable {
	tbl.ctx = ctx
	return tbl
}

// WithLogger sets the logger for subsequent queries.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) WithLogger(logger *log.Logger) CUserTable {
	tbl.logger = logger
	return tbl
}

// Logger gets the trace logger.
func (tbl CUserTable) Logger() *log.Logger {
	return tbl.logger
}

// SetLogger sets the logger for subsequent queries, returning the interface.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) SetLogger(logger *log.Logger) sqlgen2.Table {
	tbl.logger = logger
	return tbl
}

// Ctx gets the current request context.
func (tbl CUserTable) Ctx() context.Context {
	return tbl.ctx
}

// Dialect gets the database dialect.
func (tbl CUserTable) Dialect() schema.Dialect {
	return tbl.dialect
}

// Wrapper gets the user-defined wrapper.
func (tbl CUserTable) Wrapper() interface{} {
	return tbl.wrapper
}

// SetWrapper sets the user-defined wrapper.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) SetWrapper(wrapper interface{}) sqlgen2.Table {
	tbl.wrapper = wrapper
	return tbl
}

// Name gets the table name.
func (tbl CUserTable) Name() sqlgen2.TableName {
	return tbl.name
}

// DB gets the wrapped database handle, provided this is not within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl CUserTable) DB() *sql.DB {
	return tbl.db.(*sql.DB)
}

// Tx gets the wrapped transaction handle, provided this is within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl CUserTable) Tx() *sql.Tx {
	return tbl.db.(*sql.Tx)
}

// IsTx tests whether this is within a transaction.
func (tbl CUserTable) IsTx() bool {
	_, ok := tbl.db.(*sql.Tx)
	return ok
}

// Begin starts a transaction. The default isolation level is dependent on the driver.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) BeginTx(opts *sql.TxOptions) (CUserTable, error) {
	d := tbl.db.(*sql.DB)
	var err error
	tbl.db, err = d.BeginTx(tbl.ctx, opts)
	return tbl, err
}

// Using returns a modified Table using the transaction supplied. This is needed
// when making multiple queries across several tables within a single transaction.
// The result is a modified copy of the table; the original is unchanged.
func (tbl CUserTable) Using(tx *sql.Tx) CUserTable {
	tbl.db = tx
	return tbl
}

func (tbl CUserTable) logQuery(query string, args ...interface{}) {
	sqlgen2.LogQuery(tbl.logger, query, args...)
}


//--------------------------------------------------------------------------------

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// It returns the number of rows affected (of the database drive supports this).
func (tbl CUserTable) Exec(query string, args ...interface{}) (int64, error) {
	tbl.logQuery(query, args...)
	res, err := tbl.db.ExecContext(tbl.ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//--------------------------------------------------------------------------------

// Query is the low-level access method for Users.
// Note that this applies ReplaceTableName to the query string.
func (tbl CUserTable) Query(query string, args ...interface{}) ([]*User, error) {
	query = tbl.ReplaceTableName(query)
	return tbl.doQuery(false, query, args...)
}

// QueryOne is the low-level access method for one User.
// Note that this applies ReplaceTableName to the query string.
// If the query selected many rows, only the first is returned; the rest are discarded.
// If not found, *User will be nil.
func (tbl CUserTable) QueryOne(query string, args ...interface{}) (*User, error) {
	query = tbl.ReplaceTableName(query)
	return tbl.doQueryOne(query, args...)
}

func (tbl CUserTable) doQueryOne(query string, args ...interface{}) (*User, error) {
	list, err := tbl.doQuery(true, query, args...)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return list[0], nil
}

func (tbl CUserTable) doQuery(firstOnly bool, query string, args ...interface{}) ([]*User, error) {
	tbl.logQuery(query, args...)
	rows, err := tbl.db.QueryContext(tbl.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCUsers(rows, firstOnly)
}

// QueryOneNullString is a low-level access method for one string. This can be used for function queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// Note that this applies ReplaceTableName to the query string.
func (tbl CUserTable) QueryOneNullString(query string, args ...interface{}) (sql.NullString, error) {
	var result sql.NullString
	query = tbl.ReplaceTableName(query)
	tbl.logQuery(query, args...)
	rows, err := tbl.db.QueryContext(tbl.ctx, query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result)
		if err == sql.ErrNoRows {
			err = nil // not needed; result will be invalid
		}
	}
	return result, err
}

// QueryOneNullInt64 is a low-level access method for one int64. This can be used for 'COUNT(1)' queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// Note that this applies ReplaceTableName to the query string.
func (tbl CUserTable) QueryOneNullInt64(query string, args ...interface{}) (sql.NullInt64, error) {
	var result sql.NullInt64
	query = tbl.ReplaceTableName(query)
	tbl.logQuery(query, args...)
	rows, err := tbl.db.QueryContext(tbl.ctx, query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result)
		if err == sql.ErrNoRows {
			err = nil // not needed; result will be invalid
		}
	}
	return result, err
}

// QueryOneNullFloat64 is a low-level access method for one float64. This can be used for 'AVG(...)' queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// Note that this applies ReplaceTableName to the query string.
func (tbl CUserTable) QueryOneNullFloat64(query string, args ...interface{}) (sql.NullFloat64, error) {
	var result sql.NullFloat64
	query = tbl.ReplaceTableName(query)
	tbl.logQuery(query, args...)
	rows, err := tbl.db.QueryContext(tbl.ctx, query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result)
		if err == sql.ErrNoRows {
			err = nil // not needed; result will be invalid
		}
	}
	return result, err
}

// ReplaceTableName replaces all occurrences of "{TABLE}" with the table's name.
func (tbl CUserTable) ReplaceTableName(query string) string {
	return strings.Replace(query, "{TABLE}", tbl.name.String(), -1)
}

//--------------------------------------------------------------------------------

// Insert adds new records for the Users. The Users have their primary key fields
// set to the new record identifiers.
// The User.PreInsert(Execer) method will be called, if it exists.
func (tbl CUserTable) Insert(vv ...*User) error {
	var params string
	switch tbl.dialect {
	case schema.Postgres:
		params = sCUserDataColumnParamsPostgres
	default:
		params = sCUserDataColumnParamsSimple
	}

	query := fmt.Sprintf(sqlInsertCUser, tbl.name, params)
	st, err := tbl.db.PrepareContext(tbl.ctx, query)
	if err != nil {
		return err
	}
	defer st.Close()

	for _, v := range vv {
		var iv interface{} = v
		if hook, ok := iv.(sqlgen2.CanPreInsert); ok {
			hook.PreInsert()
		}

		fields, err := sliceCUserWithoutPk(v)
		if err != nil {
			return err
		}

		tbl.logQuery(query, fields...)
		res, err := st.Exec(fields...)
		if err != nil {
			return err
		}

		v.Uid, err = res.LastInsertId()
		if err != nil {
			return err
		}
	}

	return nil
}

const sqlInsertCUser = `
INSERT INTO %s (
	login,
	emailaddress,
	avatar,
	role,
	active,
	admin,
	fave,
	lastupdated,
	token,
	secret
) VALUES (%s)
`

const sCUserDataColumnParamsSimple = "?,?,?,?,?,?,?,?,?,?"

const sCUserDataColumnParamsPostgres = "$1,$2,$3,$4,$5,$6,$7,$8,$9,$10"

func sliceCUserWithoutPk(v *User) ([]interface{}, error) {

	v7, err := json.Marshal(&v.Fave)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		v.Login,
		v.EmailAddress,
		v.Avatar,
		v.Role,
		v.Active,
		v.Admin,
		v7,
		v.LastUpdated,
		v.token,
		v.secret,

	}, nil
}

// scanCUsers reads table records into a slice of values.
func scanCUsers(rows *sql.Rows, firstOnly bool) ([]*User, error) {
	var err error
	var vv []*User

	for rows.Next() {
		var v0 int64
		var v1 string
		var v2 string
		var v3 sql.NullString
		var v4 *Role
		var v5 bool
		var v6 bool
		var v7 []byte
		var v8 int64
		var v9 string
		var v10 string

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
			&v10,
		)
		if err != nil {
			return vv, err
		}

		v := &User{}
		v.Uid = v0
		v.Login = v1
		v.EmailAddress = v2
		if v3.Valid {
			a := v3.String
			v.Avatar = &a
		}
		v.Role = v4
		v.Active = v5
		v.Admin = v6
		err = json.Unmarshal(v7, &v.Fave)
		if err != nil {
			return nil, err
		}
		v.LastUpdated = v8
		v.token = v9
		v.secret = v10

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
