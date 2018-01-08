package schema

import (
	"testing"
	. "github.com/rickb777/sqlgen2/sqlgen/parse"
)

func TestDistinctTypes(t *testing.T) {
	i64 := Type{"", "", "int64", false, Int64}
	boo := Type{"", "", "bool", false, Bool}
	cat := Type{"", "", "Category", false, Int32}
	str := Type{"", "", "string", false, String}
	spt := Type{"", "", "string", true, String}
	ipt := Type{"", "", "int32", true, Int32}
	upt := Type{"", "", "uint32", true, Uint32}
	fpt := Type{"", "", "float32", true, Float32}
	sli := Type{"", "", "[]string", false, Slice}
	bgi := Type{"math/big", "big", "Int", false, Struct}
	bys := Type{"", "", "[]byte", false, Slice}
	tim := Type{"time", "time", "Time", false, Struct}

	id := &Field{Node{"Id", i64, nil}, "id", ENCNONE, Tag{Primary: true, Auto: true}}
	category := &Field{Node{"Cat", cat, nil}, "cat", ENCNONE, Tag{Index: "catIdx"}}
	name := &Field{Node{"Name", str, nil}, "username", ENCNONE, Tag{Size: 2048, Name: "username", Unique: "nameIdx"}}
	active := &Field{Node{"Active", boo, nil}, "active", ENCNONE, Tag{}}
	qual := &Field{Node{"Qual", spt, nil}, "qual", ENCNONE, Tag{}}
	diff := &Field{Node{"Diff", ipt, nil}, "diff", ENCNONE, Tag{}}
	age := &Field{Node{"Age", upt, nil}, "age", ENCNONE, Tag{}}
	bmi := &Field{Node{"Bmi", fpt, nil}, "bmi", ENCNONE, Tag{}}
	labels := &Field{Node{"Labels", sli, nil}, "labels", ENCJSON, Tag{Encode: "json"}}
	fave := &Field{Node{"Fave", bgi, nil}, "fave", ENCJSON, Tag{Encode: "json"}}
	avatar := &Field{Node{"Avatar", bys, nil}, "avatar", ENCNONE, Tag{}}
	updated := &Field{Node{"Updated", tim, nil}, "updated", ENCTEXT, Tag{Size: 100, Encode: "text"}}

	//icat := &Index{"catIdx", false, FieldList{category}}
	//iname := &Index{"nameIdx", true, FieldList{name}}

	cases := []struct {
		list     FieldList
		expected TypeSet
	}{
		{FieldList{id}, NewTypeSet(i64)},
		{FieldList{id, id, id}, NewTypeSet(i64)},
		{FieldList{id, category}, NewTypeSet(i64, cat)},
		{FieldList{id,
			category,
			name,
			qual,
			diff,
			age,
			bmi,
			active,
			labels,
			fave,
			avatar,
			updated}, NewTypeSet(i64, boo, cat, str, spt, ipt, upt, fpt, bgi, sli, bys, tim)},
	}
	for _, c := range cases {
		s := c.list.DistinctTypes()
		if !NewTypeSet(s...).Equals(c.expected) {
			t.Errorf("expected %d::%+v but got %d::%+v", len(c.expected), c.expected, len(s), s)
		}
	}
}
