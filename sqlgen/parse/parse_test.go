package parse

import (
	"bytes"
	"github.com/kortschak/utter"
	"reflect"
	"strings"
	"testing"
	"fmt"
	"github.com/rickb777/sqlgen2/sqlgen/parse/exit"
	"os"
)

func TestStructWith3FieldsAndTags(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg1", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Id",
					TypeRef: NewTypeRef("", "int64", Int64),
					Tags: &Tag{Primary: true, Auto: true},
				}, {
					Name: "Number",
					TypeRef: NewTypeRef("", "int", Int),
					Tags: &Tag{},
				}, {
					Name: "Title",
					TypeRef: NewTypeRef("", "string", String),
					Tags: &Tag{},
				}, {
					Name: "Description",
					TypeRef: NewTypeRef("", "string", String),
					Tags: &Tag{},
				}, {
					Name: "Owner",
					TypeRef: NewTypeRef("", "string", String),
					Tags: &Tag{},
				},
			},
		},
		"pkg1", "Struct",
		`package pkg1

		type Struct struct {
			Id       int64 |sql:"pk: true, auto: true"|
			Number   int
			Title, Description, Owner    string // must find all three fields
		}`,
	)
}

func TestStructWith1BoolFieldAndIgnoreTag(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg2", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Flag",
					TypeRef: NewTypeRef("", "bool", Bool),
					Tags: &Tag{Skip: true},
				},
			},
		},
		"pkg2", "Struct",
		`package pkg2

		type Struct struct {
			Flag  bool     |sql:"-"|
		}`,
	)
}

func TestStructWith1SliceFieldAndJsonTag(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg3", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Labels",
					TypeRef: NewTypeRef("", "[]string", Slice),
					Tags: &Tag{Encode: "json"},
				},
			},
		},
		"pkg3", "Struct",
		`package pkg3

		type Struct struct {
			Labels   []string  |sql:"encode: json"|
		}`,
	)
}

func TestStructWith1MapFieldAndJsonTag(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg4", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Table",
					TypeRef: NewTypeRef("", "map[string]int", Map),
					Tags: &Tag{Encode: "json"},
				},
			},
		},
		"pkg4", "Struct",
		`package pkg4

		type Struct struct {
			Table    map[string]int  |sql:"encode: json"|
		}`,
	)
}

func TestStructWithNestedStructType(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg5", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Author",
					TypeRef: NewTypeRef("", "Author", Ptr),
					Tags: &Tag{},
					Nodes: []*Node{
						{
							Name: "Name",
							TypeRef: NewTypeRef("", "string", String),
							Tags: &Tag{},
						},
					},
				},
			},
		},
		"pkg5", "Struct",
		`package pkg5

		type Struct struct {
			Author    *Author
		}
		type Author struct {
			Name     string
		}`,
	)
}

func TestStructWithNestedSimpleType(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg6", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Cat",
					TypeRef: NewTypeRef("pkg6", "Category", Int32),
					Tags: &Tag{},
				},
			},
		},
		"pkg6", "Struct",
		`package pkg6

		type Category int32

		type Struct struct {
			Cat      Category
		}`,
	)
}

func TestStructWithNestedSimpleTypeInOtherPackageOrder1(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg7", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Cat",
					TypeRef: NewTypeRef("other", "Category", Int32),
					Tags: &Tag{},
				},
			},
		},
		"pkg7", "Struct",
		`package pkg7

		type Struct struct {
			Cat      other.Category
		}`,
		`package other

		type Category int32
		`,
	)
}

func TestStructWithNestedSimpleTypeInOtherPackageOrder2(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg8", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Cat",
					TypeRef: NewTypeRef("other", "Category", Int32),
					Tags: &Tag{},
				},
			},
		},
		"pkg8", "Struct",
		`package other

		type Category int32
		`,
		`package pkg8

		type Struct struct {
			Cat      other.Category
		}`,
	)
}

func TestStructWithNestingAcross2Packages(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg9", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Id",
					TypeRef: NewTypeRef("", "uint32", Uint32),
					Tags: &Tag{},
				},
				{
					Name: "Wibble",
					TypeRef: NewTypeRef("stringy", "Thingy", String),
					Tags: &Tag{},
				},
				{
					Name: "Bobble",
					TypeRef: NewTypeRef("", "string", String),
					Tags: &Tag{},
				},
			},
		},
		"pkg9", "Struct",
		`package stringy

		type Thingy string
		`,
		`package froob

		type Inner1 struct {
			Wibble stringy.Thingy
		}
		`,

		`package pkg9

		type Struct struct {
			Id uint32
			froob.Inner1
			Bobble string
		}
		`,
	)
}

func TestStructWithNestingInTheSamePackage(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg10", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Id",
					TypeRef: NewTypeRef("", "uint32", Uint32),
					Tags: &Tag{},
				},
				{
					Name: "Uid",
					TypeRef: NewTypeRef("", "uint32", Uint32),
					Tags: &Tag{},
				},
				{
					Name: "Name",
					TypeRef: NewTypeRef("pkg10", "Username", String),
					Tags: &Tag{},
				},
				{
					Name: "Wibble",
					TypeRef: NewTypeRef("pkg10", "Thingy", String),
					Tags: &Tag{},
				},
				{
					Name: "Bobble",
					TypeRef: NewTypeRef("pkg10", "Username", String),
					Tags: &Tag{},
				},
			},
		},
		"pkg10", "Struct",
		`package pkg10

		type Thingy string

		type Username string

		type User struct {
			Uid  uint32
			Name Username
		}

		type UserWithThingy struct {
			User
			Wibble Thingy
		}

		type Struct struct {
			Id uint32
			UserWithThingy
			Bobble Username
		}
		`,
	)
}

func TestStructWithNestingAcross4Packages(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg11", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Id",
					TypeRef: NewTypeRef("", "uint32", Uint32),
					Tags: &Tag{},
				},
				{
					Name: "Uid",
					TypeRef: NewTypeRef("", "uint32", Uint32),
					Tags: &Tag{},
				},
				{
					Name: "Name",
					TypeRef: NewTypeRef("userindex", "Username", String),
					Tags: &Tag{},
				},
				{
					Name: "Wibble",
					TypeRef: NewTypeRef("stringy", "Thingy", String),
					Tags: &Tag{},
				},
				{
					Name: "Bobble",
					TypeRef: NewTypeRef("userindex", "Username", String),
					Tags: &Tag{},
				},
				{
					Name: "Date",
					TypeRef: NewTypeRef("date", "PeriodOfDays", Int32),
					Tags: &Tag{},
				},
			},
		},
		"pkg11", "Struct",
		`package stringy

		type Thingy string
		`, //--------------------------

		`package userindex

		type Username string
		`, //--------------------------

		`package userindex

		type User struct {
			Uid  uint32
			Name Username
		}

		type UserWithThingy struct {
			User
			Wibble stringy.Thingy
		}
		`, //--------------------------

		`package date

		type PeriodOfDays int32

		type Date struct {
			day PeriodOfDays
		}
		`, //--------------------------

		`package pkg11

		type Struct struct {
			Id uint32
			userindex.UserWithThingy
			Bobble userindex.Username
			Date   date.Date
		}
		`,
	)
}

func TestStructWithUnexportedFields(t *testing.T) {
	doTestParseOK(t,
		&Node{
			Name: "Struct",
			TypeRef: NewTypeRef("pkg12", "Struct", Struct),
			Nodes: []*Node{
				{
					Name: "Date",
					TypeRef: NewTypeRef("date", "Date", Scanner),
					Tags: &Tag{},
				},
			},
		},
		"pkg12", "Struct",
		`package date

		type PeriodOfDays int32

		type Date struct {
			day PeriodOfDays
		}
		`, //--------------------------

		`package pkg12

		type Struct struct {
			Date   date.Date
		}
		`,
	)
}

func doTestParseOK(t *testing.T, want *Node, pkg, name string, isource ...string) {
	t.Helper()
	exit.TestableExit()
	//Debug = true

	// fix edges missing in the literal values
	for _, n0 := range want.Nodes {
		n0.Parent = want
		for _, n1 := range n0.Nodes {
			n1.Parent = n0
		}
	}

	files := make([]Source, len(isource))

	for i, s := range isource {
		// allow nested back-ticks
		source := strings.Replace(s, "|", "`", -1)
		files[i] = Source{fmt.Sprintf("issue%d.go", i), bytes.NewBufferString(source)}
	}

	err := parseAllFiles(files)
	if err != nil {
		t.Errorf("Error parsing: %s", err)
	}

	got, err := findMatchingNodes(pkg, name)
	if err != nil {
		t.Errorf("Error parsing: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		ex := utter.Sdump(want)
		ac := utter.Sdump(got)
		outputDiff(ex, pkg+"-expected.txt")
		outputDiff(ac, pkg+"-got.txt")
		t.Errorf("Wanted %s\nGot %s", ex, ac)
	}
}

//-------------------------------------------------------------------------------------------------

func outputDiff(a, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	f.WriteString(a)
	f.WriteString("\n")
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
