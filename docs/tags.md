# Tag In More Detail

Tags are optional annotations on the fields of the structs used to define your tables in Go. Append them after the field types like this.

```go
type User struct {
    ID      int64             `sql:"pk: true, auto: true"`
    Login   string            `sql:"unique: loginIdx"`
    Email   string            `sql:"index: emailIdx"`
    Age     uint              `sql:"check: age>=18"`
    Prefs   map[string]string `sql:"encode: json"`
}
```

The double-quoted string is actually treated as a snippet of YAML, so it needs some whitespace or it won't be accepted.


## Tags

| Tag      | Value         | Purpose                                                           |
| -------- | ------------- | ----------------------------------------------------------------- |
| pk       | true or false | the column is the primary key                                     |
| auto     | true or false | the column is auto-incrementing (ignored if not using MySQL)      |
| nk       | true or false | the column is the natural key; a unique index will be created     |
| prefixed | true or false | the column name is made unique using a computed prefix            |
| name     | string        | the column name (default is the field's name)                     |
| type     | string        | overrides the column type explicitly                              |
| default  | string        | sets a default value in the DDL (only for non-nullable columns)   |
| size     | integer       | sets the storage size for the column                              |
| encode   | string        | encodes as "json", "text" or using the "driver"                   |
| index    | string        | the column has an index                                           |
| unique   | string        | the column has a unique index                                     |
| check    | string        | apply the expression as a check constraint on writes to the table |
| fk       | string        | the column is a foreign key of some parent table                  |
| onupdate | string        | one of "cascade", "delete", "restrict", "set default", "set null" |
| ondelete | string        | one of "cascade", "delete", "restrict", "set default", "set null" |

For string tags, you need to surround the value in single quotes if it contains any spaces. For example, `sql:"check: 'age >= 18'"`. You don't need quotes if there aren't any spaces.

Marking a field to be a natural key is the same as adding a single-column unique index. A common pattern is to define a natural key based on the business properties of the field (for example all email addresses in a table are unique and provide normative identity), and also provide a surrogate key for the primary key using an int64 (or similar) as the internal identity. Natural keys often make inferior foreign keys; surrogate keys are usually better.


## Indexes

If your columns need indexes, `sqlgen` includes extra code for CRUD operations based on the indexed columns as well as on the primary key. This example shows a primary key column `Uid`, an indexed column `Name `, and a uniquely indexed column `EmailAddress`.

```Go
type User struct {
	Uid          int64    `sql:"pk: true, auto: true"`
	Name         string   `sql:"index: user_name"`
	EmailAddress string   `sql:"unique: user_email"`
    // ... other fields not showm
}
```

## Encode

* "json" - the struct field will be marshalled using Go's JSON marshaller and stored as a string of JSON.
* "text" - the struct field will be marshalled using a TextMarshaler and stored as a string. You have to provide the TextMarshaller.
* "driver" - the struct field will be marshalled by Go's SQL driver using interfaces implemented by you: `sql.Scanner` and `driver.Valuer`. The column type is whatever you want; your `Scanner` and `Valuer` functions must support the type you choose.


### JSON Encoding

Some types in your struct may not have native equivalents in your database, such as `[]string`. These values can be marshalled and stored as JSON in the database instead.

```diff
type User struct {
    Uid    int64  `sql:"pk: true"`
    Name   string
    Email  string
+   Label  []string `sql:"encode: json"
}
```

You don't have to do anything else; the jSON marshal and unmarshal code will be added automatically.


### Driver Encoding

If you have a field type that implements `sql.Scanner`, `driver.Valuer` and you want to use this explicitly as the column encoding to string data, specify `encode: driver`, e.g.

```
type User struct {
    Uid    int64    `sql:"pk: true"`
    Stuff  MyStruct `sql:"encode: driver"`
}
```

You don't always have to do this because your types (such as `MyStruct`) are inspected and will normally be auto-detected. However, for struct types that contain other fields, ambiguity arises. You can use this setting to resolve the ambiguity; otherwise the internal fields will be treated as table columns.

In the example above, `Stuff` is treated as a single database column and `MyStruct` must provide `Scan` and `Value` methods. No fields within `MyStruct` are inspected; they are ignored instead.


### Text Encoding

If you have a field type that implements `encoding.MarshalText`, `encoding.UnmarshalText` and you want to use this as the column encoding to string data, specify `encode: text`, e.g.

```
type User struct {
    Uid    int64    `sql:"pk: true"`
    Thing  MyStruct `sql:"encode: text"`
}
```

There is no auto-detection for these encoding interfaces. Use this setting as and when you need it. The compiler will tell you if you got it wrong (e.g. the `MyStruct` type was missing one of the required methods).


## Check Constraints

Use a tag such as `sql:"check: age>=18"`when you want the value of a column to be constrained by the expression you specify.

More info: https://www.postgresql.org/docs/9.6/static/ddl-constraints.html#DDL-CONSTRAINTS-CHECK-CONSTRAINTS


## Foreign Key Constraints

When `fk` is specified, it must contain a table name and the primary key in that table. It's in the form `table.pk`. Example:

```Go
type User struct {
	Uid          int64    `sql:"pk: true, auto: true"`
	// ...
	AddressId    *int64   `sql:"fk: addresses.id, onupdate: restrict, ondelete: restrict"`
	// ... other fields not shown
}
```

This describes a `users` table that refers to an `addresses` table using a foreign key, which is optional in this case (i.e. it's a pointer). The primary key of table `addresses` is a column called `id` and so the foreign key is denoted using `fk: addresses.id`.

The `onupdate` and `ondelete` options cause `... ON UPDATE ...` and `... ON DELETE ...` clauses being to be added to the resulting SQL.


## See Also

 * back to [**the tutorial**](tutorial.md)
