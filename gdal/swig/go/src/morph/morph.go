package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

type MergeVisitor struct {
	l *LoadVisitor
	r map[string]string
}

func NewMergeVisitor(l *LoadVisitor) *MergeVisitor {
	r := map[string]string{
		"ReadDir":           "wrap_ReadDir",
		"SetErrorHandler":   "wrap_SetErrorHandler",
		"Open":              "wrap_Open",
		"ApplyGeoTransform": "wrap_ApplyGeoTransform",
		"InvGeoTransform":   "wrap_InvGeoTransform",
	}

	return &MergeVisitor{l, r}
}

func (v *MergeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.TypeSpec:
		switch it := t.Type.(type) {
		case *ast.InterfaceType:
			methods, ok := v.l.Tree[t.Name.Name]
			if !ok {
				break
			}
			log.Printf("Merging interface %s", t.Name.Name)
			it.Methods.List = append(it.Methods.List, methods.List...)
		}
	case *ast.FuncDecl:
		new, ok := v.r[t.Name.Name]
		if !ok {
			break
		}
		log.Printf("Renaming %v to %v", t.Name.Name, new)
		t.Name = ast.NewIdent(new)
	}

	return v
}

type LoadVisitor struct {
	Tree map[string]*ast.FieldList
}

func NewLoadVisitor() *LoadVisitor {
	return &LoadVisitor{make(map[string]*ast.FieldList)}
}

func (v *LoadVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.TypeSpec:
		switch it := t.Type.(type) {
		case *ast.InterfaceType:
			log.Printf("Loading interface %s", t.Name.Name)
			v.Tree[t.Name.Name] = it.Methods
		}
	}

	return v
}

func load() (v *LoadVisitor) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "merge/merge.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	v = NewLoadVisitor()
	ast.Walk(v, file)
	return v
}

func main() {
	l := load()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "build/gdal.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	v := NewMergeVisitor(l)
	ast.Walk(v, file)

	out, err := os.Create("src/osgeo/gdal.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	printer.Fprint(out, fset, file)
}
