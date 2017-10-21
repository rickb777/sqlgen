package schema

type DialectId int

const (
	SQLITE DialectId = iota
	POSTGRES
	MYSQL
)

var Dialects = map[string]DialectId{
	"sqlite":   SQLITE,
	"postgres": POSTGRES,
	"mysql":    MYSQL,
}

type Dialect interface {
	Id() DialectId
	Table(*Table) string
	Index(*Table, *Index) string
	Insert(*Table) string
	Update(*Table, []*Field) string
	Delete(*Table, []*Field) string
	Param(int) string
	Params(int, int) []string
	ColumnParams(t *Table, withAuto bool) string
	Token(int) string
}

func New(dialect DialectId) Dialect {
	switch dialect {
	case POSTGRES:
		return newPostgres()
	case MYSQL:
		return newMySQL()
	default:
		return newSQLite()
	}
}

func (f *Field) Column(dialect DialectId) string {
	switch dialect {
	case MYSQL:
		return mysqlColumn(f)
	case POSTGRES:
		return postgresColumn(f)
	case SQLITE:
		return sqliteColumn(f)
	default:
		return ""
	}
}

