package code

import (
	"fmt"
	"io"
	. "strings"

	"bitbucket.org/pkg/inflect"
	"github.com/rickb777/sqlgen2/schema"
)

const sectionBreak = "\n//--------------------------------------------------------------------------------"

type ConstView struct {
	Name string
	Body interface{}
}

//const constString = "\nconst %s = %s\n"
const constStringQ = "\nconst %s = %q\n"
const constStringWithTicks = "\nconst %s = `\n%s\n`\n"

// WritePackage writes the Go package header to
// writer w with the given package name.
func WritePackage(w io.Writer, name string) {
	fmt.Fprintf(w, sPackage, name)
}

// writeSchema writes SQL statements to CREATE, INSERT,
// UPDATE and DELETE values from Table t.
func WriteSchema(w io.Writer, view View, table *schema.Table) {

	fmt.Fprintln(w, sectionBreak)

	tableName := view.Prefix + table.Type

	for _, did := range schema.AllDialectIds {
		d := schema.New(did)
		ds := did.String()

		fmt.Fprintf(w, constStringWithTicks,
			identifier("sqlCreate", tableName, "Table"+ds),
			"CREATE TABLE %s%s%s ("+d.Table(table, did)+"\n)"+d.CreateTableSettings())
	}

	if table.HasPrimaryKey() {
		fmt.Fprintln(w, sectionBreak)

		fmt.Fprintf(w, constStringWithTicks,
			identifier("sqlDelete", tableName, "ByPkPostgres"), schema.New(schema.Postgres).Delete(table, []*schema.Field{table.Primary}))

		fmt.Fprintf(w, constStringWithTicks,
			identifier("sqlDelete", tableName, "ByPkSimple"), schema.New(schema.Sqlite).Delete(table, []*schema.Field{table.Primary}))

	}

	for _, did := range schema.AllDialectIds {
		d := schema.New(did)
		ds := did.String()

		for _, ix := range table.Index {

			fmt.Fprintln(w, sectionBreak)

			fmt.Fprintf(w, constStringWithTicks,
				identifier("sqlCreate"+view.Prefix, ix.Name, "Index"+ds), d.Index(table, ix))

			if ix.Unique {
				fmt.Fprintf(w, constStringWithTicks,
					identifier("sqlUpdate"+view.Prefix, ix.Name, ds), d.Update(table, ix.Fields))

				fmt.Fprintf(w, constStringWithTicks,
					identifier("sqlDelete"+view.Prefix, ix.Name, ds), d.Delete(table, ix.Fields))
			}
		}
	}

	fmt.Fprintln(w, sectionBreak)

	fmt.Fprintf(w, "\nconst %s = %d\n",
		identifier("Num", tableName, "Columns"), table.NumColumnNames(true))

	fmt.Fprintf(w, "\nconst %s = %d\n",
		identifier("Num", tableName, "DataColumns"), table.NumColumnNames(false))

	if table.HasPrimaryKey() {
		fmt.Fprintf(w, constStringQ,
			identifier("", tableName, "Pk"), table.Primary.Name)
	}

	fmt.Fprintf(w, constStringQ,
		identifier("", tableName, "DataColumnNames"), Join(table.ColumnNames(false), ", "))

	fmt.Fprintln(w, sectionBreak)
}

func WriteCreateTableFunc(w io.Writer, view View, table *schema.Table) {
	must(tCreateTable.Execute(w, view))
}

func WriteCreateIndexFunc(w io.Writer, view View, table *schema.Table) {
	for _, ix := range table.Index {
		view.Body1 = append(view.Body1, inflect.Camelize(ix.Name))
	}
	must(tCreateIndex.Execute(w, view))
}

func identifier(prefix, id, suffix string) string {
	return prefix + inflect.Camelize(id) + suffix
}