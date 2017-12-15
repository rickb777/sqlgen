package main

import (
	"fmt"
	"go/types"
	"strings"
	. "github.com/acsellers/inflections"
	. "github.com/rickb777/sqlgen2/schema"
	"github.com/rickb777/sqlgen2/sqlgen/parse"
	"github.com/rickb777/sqlgen2/sqlgen/parse/exit"
	"github.com/rickb777/sqlgen2/sqlgen/output"
	"github.com/kortschak/utter"
)

type context struct {
	pkgStore         parse.PackageStore
	indices          map[string]*Index
	table            *TableDescription
	unexportedFields []string
}

func load(pkgStore parse.PackageStore, pkg, name string) (*TableDescription, error) {
	table := new(TableDescription)

	// local map of indexes, used for quick lookups and de-duping.
	indices := map[string]*Index{}

	table.Type = name
	table.Name = Pluralize(Underscore(table.Type))

	str, tags := pkgStore.Find(pkg, name)
	ctx := &context{pkgStore, indices, table, nil}
	ctx.examineStruct(str, pkg, name, tags, nil)

	if len(ctx.unexportedFields) > 0 {
		output.Info("Warning: %s.%s contains unexported fields %s"+
			" (perhaps annotate with `sql:\"-\"`).\n", pkg, name, strings.Join(ctx.unexportedFields, ", "))
	}

	for _, idx := range ctx.indices {
		table.Index = append(table.Index, idx)
	}

	checkNoConflictingNames(pkg, name, table)

	return table, nil
}

//-------------------------------------------------------------------------------------------------

func (ctx *context) examineStruct(str *types.Struct, pkg, name string, tags map[string]parse.Tag, parent *Node) {
	parse.DevInfo("examineStruct %s %s\n  tags %v\n", pkg, name, tags)
	if str.NumFields() == 0 {
		exit.Fail(1, "%s.%s: empty structs are not supported (was there a parser warning?).\n", pkg, name)
	}

	for j := 0; j < str.NumFields(); j++ {
		tField := str.Field(j)
		parse.DevInfo("    f%d: name:%-10s pkg:%s type:%-25s f:%v, e:%v, a:%v\n", j,
			tField.Name(), tField.Pkg().Name(), tField.Type(), tField.IsField(), tField.Exported(), tField.Anonymous())

		if !tField.Exported() {
			if tag, exists := tags[tField.Name()]; !exists || (exists && !tag.Skip) {
				ctx.unexportedFields = append(ctx.unexportedFields, tField.Name())
			}
		}

		if tField.Anonymous() {
			ctx.convertEmbeddedNodeToFields(tField, pkg, parent)

		} else {
			ctx.convertLeafNodeToField(tField, pkg, tags, parent)
		}
	}
}

//-------------------------------------------------------------------------------------------------

func (ctx *context) convertEmbeddedNodeToFields(leaf *types.Var, pkg string, parent *Node) {
	name := leaf.Name()
	parse.DevInfo("convertEmbeddedNodeToFields %s %s\n", pkg, name)
	str, tags := ctx.pkgStore.Find(pkg, name)
	if str == nil {
		nm, ok := leaf.Type().(*types.Named)
		if !ok {
			exit.Fail(5, "Unable to find %s.%s\n", pkg, name)
		}
		pkg = nm.Obj().Pkg().Name()
		str = nm.Underlying().(*types.Struct)
		tags = make(map[string]parse.Tag)
		addStructTags(tags, str)
		parse.DevInfo(" - found in other package %v %v\n", leaf.Type(), str)
	}
	node := &Node{Name: name, Parent: parent}
	ctx.examineStruct(str, pkg, name, tags, node)
}

//-------------------------------------------------------------------------------------------------

func (ctx *context) convertLeafNodeToField(leaf *types.Var, pkg string, tags map[string]parse.Tag, parent *Node) {
	tag := tags[leaf.Name()]

	if tag.Skip {
		return
	}

	field := &Field{}
	field.Tags = tag
	field.Encode = mapTagToEncoding[tag.Encode]

	// only recurse into the node's fields if the leaf isn't encoded
	var ok bool
	field.Node, ok = ctx.convertLeafNodeToNode(leaf, pkg, tags, parent, field.Encode == ENCNONE)
	if !ok {
		return
	}

	// Lookup the SQL column type
	field.SqlType = BLOB
	underlying := leaf.Type().Underlying()
	switch u := underlying.(type) {
	case *types.Basic:
		field.SqlType = mapKindToSqlType[u.Kind()]
		field.Type.Base = parse.Kind(u.Kind())

	case *types.Slice:
		field.Type.Base = parse.Slice
	}

	if tag.Encode == "json" {
		field.SqlType = JSON
	}

	if tag.Primary {
		if ctx.table.Primary != nil {
			exit.Fail(1, "%s, %s: compound primary keys are not supported.\n",
				ctx.table.Primary.Type.Name, field.Type.Name)
		}
		ctx.table.Primary = field
	}

	if tag.Index != "" {
		index, ok := ctx.indices[tag.Index]
		if !ok {
			index = &Index{
				Name: tag.Index,
			}
			ctx.indices[index.Name] = index
		}
		index.Fields = append(index.Fields, field)
	}

	if tag.Unique != "" {
		index, ok := ctx.indices[tag.Unique]
		if !ok {
			index = &Index{
				Name: tag.Unique,
			}
			ctx.indices[index.Name] = index
		}
		index.Fields = append(index.Fields, field)
		index.Unique = true
	}

	if tag.Type != "" {
		t, ok := mapStringToSqlType[tag.Type]
		if ok {
			field.SqlType = t
		}
	}

	prefix := ""
	if tag.Prefixed {
		prefix = Underscore(field.JoinParts(1, "_")) + "_"
	}

	if tag.Name != "" {
		field.SqlName = prefix + strings.ToLower(tag.Name)
	} else {
		field.SqlName = prefix + strings.ToLower(field.Name)
	}

	ctx.table.Fields = append(ctx.table.Fields, field)
}

//-------------------------------------------------------------------------------------------------

func (ctx *context) convertLeafNodeToNode(leaf *types.Var, pkg string, tags map[string]parse.Tag, parent *Node, canRecurse bool) (Node, bool) {
	node := Node{Name: leaf.Name(), Parent: parent}
	tp := Type{}

	lt := leaf.Type()
	//isPtr := false

	switch t := lt.(type) {
	case *types.Pointer:
		lt = t.Elem()
		tp.Base = parse.Ptr
		//isPtr = true
	}

	switch t := lt.(type) {
	case *types.Basic:
		tp.Name = t.Name()
		//case *types.Struct:
		//	tp.Name = t.String()
	case *types.Named:
		tObj := t.Obj()
		if tObj.Pkg().Name() != pkg {
			tp.PkgPath = tObj.Pkg().Path()
			tp.PkgName = tObj.Pkg().Name()
		}
		tp.Name = tObj.Name()

		if str, ok := t.Underlying().(*types.Struct); ok {
			tp.Base = parse.Struct
			if canRecurse {
				addStructTags(tags, str)
				ctx.examineStruct(str, pkg, leaf.Name(), tags, &node)
				return node, false
			}
		}
	case *types.Array:
		tp.Name = t.String()
	case *types.Slice:
		switch el := t.Elem().(type) {
		case *types.Basic:
			tp.Name = t.String()
		case *types.Named:
			tnObj := el.Obj()
			parse.DevInfo("slice pkgname:%s pkgpath:%s name:%s\n", tnObj.Pkg().Name(), tnObj.Pkg().Path(), tnObj.Name())
			if tnObj.Pkg().Name() != pkg {
				tp.PkgPath = tnObj.Pkg().Path()
				tp.PkgName = tnObj.Pkg().Name()
			}
			tp.Name = tnObj.Name()
		}
	default:
		panic(fmt.Sprintf("%#v", lt))
	}

	node.Type = tp
	return node, true
}

func addStructTags(tags map[string]parse.Tag, str *types.Struct) {
	for i := 0; i < str.NumFields(); i++ {
		ts := str.Tag(i)
		tag, err := parse.ParseTag(ts)
		if err != nil {
			exit.Fail(2, "%s contains unparseable tag %q (%s)", str.String(), ts, err)
		}
		tags[str.Field(i).Name()] = *tag
	}
}

func checkNoConflictingNames(pkg, name string, table *TableDescription) {
	names := make(map[string]struct{})
	var duplicates []string

	for _, field := range table.Fields {
		name := strings.ToLower(field.SqlName)
		_, exists := names[name]
		if exists {
			duplicates = append(duplicates, name)
		}
		names[name] = struct{}{}
	}

	if len(duplicates) > 0 {
		parse.DevInfo("checkNoConflictingNames %s %s %+v\n", pkg, name, utter.Sdump(table))
		exit.Fail(1, "%s.%s: found conflicting SQL column names: %s.\nPlease set the names on these fields explicitly using tags.\n",
			pkg, name, strings.Join(duplicates, ", "))
	}
}