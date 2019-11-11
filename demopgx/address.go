package demopgx

// This example demonstrates
//   * JSON encoding for a slice of strings
//   * the use of pointers for optional items
//   * metadata is written to/read from JSON file

//go:generate sqlgen -pgx -type demopgx.Address -o address_sql.go -all -v address.go

type AddressFields struct {
	Lines    []string `sql:"encode: json"`
	Town     *string  `sql:"size: 80, index: townIdx"`
	Postcode string   `sql:"size: 20, index: postcodeIdx"`
	UPRN     string   `sql:"nk: true, size: 20"`
}

type Address struct {
	Id int64 `sql:"pk: true, auto: true"`
	AddressFields
}
