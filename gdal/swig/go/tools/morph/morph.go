package main

import (
	"encoding/json"
	"flag"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type ActionVisitor struct {
	l *LoadVisitor
	f Functions
}

func NewActionVisitor(l *LoadVisitor, f Functions) *ActionVisitor {
	return &ActionVisitor{l, f}
}

func (v *ActionVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.TypeSpec:
		switch it := t.Type.(type) {
		case *ast.InterfaceType:
			if v.f != nil { // Delete methods.
				w := 0 // write index
			loop:
				for _, field := range it.Methods.List {
					fname := field.Names[0].Name
					//log.Printf("Method %v.%v()", t.Name.Name, fname)
					if action, ok := v.f[fname]; ok == true {
						for _, iface := range action.Delete {
							if iface == t.Name.Name {
								log.Printf("Removing interface method %v.%v()", iface, fname)
								continue loop
							}
						}
					}
					it.Methods.List[w] = field
					w++
				}
				it.Methods.List = it.Methods.List[:w]
			}

			if v.l != nil { // Merge methods.
				methods, ok := v.l.Tree[t.Name.Name]
				if !ok {
					break
				}
				log.Printf("Merging %d methods from interface %s", len(methods.List), t.Name.Name)
				it.Methods.List = append(it.Methods.List, methods.List...)
			}
		}
	case *ast.FuncDecl:
		if v.f != nil { // Rename functions and methods.
			fname := t.Name.Name
			if action, ok := v.f[fname]; ok == true {
				if action.Rename == nil {
					break
				}
				r := action.Rename

				if t.Recv != nil && t.Recv.List != nil && len(t.Recv.List) == 1 { // Is it a method?
					switch st := t.Recv.List[0].Type.(type) {
					case *ast.Ident:
						if rename, ok := r.Method[st.Name]; ok == true {
							new := rename.Rename(t.Name.Name)
							log.Printf("Renaming %v.%v() to %v.%v()", st.Name, t.Name.Name, st.Name, new)
							t.Name = ast.NewIdent(new)
						}
					}
				} else { // It's a function.
					if r.Function != nil {
						new := r.Function.Rename(t.Name.Name)
						log.Printf("Renaming %v() to %v()", t.Name.Name, new)
						t.Name = ast.NewIdent(new)
					}
				}
			}
		}
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

func loadMerge(fname string) (v *LoadVisitor, err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		return
	}

	v = NewLoadVisitor()
	ast.Walk(v, file)
	return v, nil
}

type Functions map[string]*Actions

type Actions struct {
	Delete []string
	Rename *Renames
}

type Renames struct {
	Function *Rename
	Method   map[string]*Rename
}

type Rename struct {
	Prefix   string
	Unexport bool
}

func (r *Rename) Rename(name string) string {
	name = r.Prefix + name
	if r.Unexport {
		r, size := utf8.DecodeRuneInString(name)
		name = strings.ToLower(string(r)) + name[size:]
	}
	return name
}

func loadFunctions(fname string) (Functions, error) {
	var funcs Functions

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err := dec.Decode(&funcs); err != nil {
		return nil, err
	}

	return funcs, nil
}

var fmods = flag.String("mods", "", "json file defining modifications to functions and methods")
var fmerge = flag.String("merge", "", "go source file defining interfaces to merge")
var fin = flag.String("in", "", "input go source file (defaults to STDIN)")
var fout = flag.String("out", "", "output go source file (defaults to STDOUT)")

func main() {
	flag.Parse()

	var funcs Functions
	var err error
	if len(*fmods) > 0 {
		funcs, err = loadFunctions(*fmods)
		if err != nil {
			log.Fatal(err)
		}
	}

	var merge *LoadVisitor
	if len(*fmerge) > 0 {
		merge, err = loadMerge(*fmerge)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the input from STDIN by default.
	if len(*fin) == 0 {
		fh, e := ioutil.TempFile("", "input")
		if e != nil {
			log.Fatal(e)
		}

		_, err = io.Copy(fh, os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		*fin = fh.Name()
		defer func() {
			fh.Close()
			os.Remove(*fin)
		}()
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *fin, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	if merge != nil || len(funcs) > 0 {
		v := NewActionVisitor(merge, funcs)
		ast.Walk(v, file)
	}

	// Get a Writer for output, defaulting to STDOUT.
	var out io.WriteCloser
	if len(*fout) > 0 {
		out, err = os.Create(*fout)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	printer.Fprint(out, fset, file)
}
