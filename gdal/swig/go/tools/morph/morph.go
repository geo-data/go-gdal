package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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

type funcType int

const (
	swigFunc funcType = iota
	overloadFunc
	apiFunc
)

func swigFuncType(fname string) funcType {
	if strings.Contains(fname, "__SWIG_") {
		return overloadFunc
	} else if strings.Contains(strings.ToLower(fname), "swig") {
		return swigFunc
	}
	return apiFunc
}

func appendErrResult(f *ast.FuncType) bool {
	if f.Results == nil {
		f.Results = &ast.FieldList{}
	}

	if f.Results.List == nil {
		f.Results.List = []*ast.Field{}
	}

	if len(f.Results.List) > 0 && f.Results.List[0].Names == nil { // it's an anonymous field.
		f.Results.List = append(f.Results.List, anonErrField)
		return true
	} else {
		f.Results.List = append(f.Results.List, errField)
		return false
	}
}

func (v *ActionVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.TypeSpec:
		switch it := t.Type.(type) {
		case *ast.InterfaceType:
			if v.f != nil { // Delete or rename interface methods.
				iname := t.Name.Name // interface name
				w := 0               // write index
			loop:
				for _, field := range it.Methods.List {
					fname := field.Names[0].Name
					//log.Printf("Method %v.%v()", t.Name.Name, fname)

					if action, ok := v.f[fname]; ok == true {
						// Delete interfaces.
						if action.Delete != nil {
							for _, iface := range action.Delete {
								if iface == iname {
									log.Printf("Removing interface method %v.%v()", iface, fname)
									continue loop
								}
							}
						}

						// Rename interfaces.
						if action.Rename != nil && action.Rename.Interface != nil {
							r := action.Rename
							if rename, ok := r.Interface[iname]; ok {
								new := rename.Rename(fname)
								log.Printf("Renaming %v.%v() to %v.%v()", iname, fname, iname, new)
								field.Names[0] = ast.NewIdent(new)
							}
						}
					}

					_, skip := skipErrors[fname]
					_, isPanic := panicErrors[fname]
					if !skip && !isPanic && swigFuncType(fname) != swigFunc {
						log.Printf("Adding error return value to %v.%v()", t.Name.Name, fname)
						sig := field.Type.(*ast.FuncType)
						appendErrResult(sig)
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
		fname := t.Name.Name

		// Add error traps to all non swig functions and methods.
		_, skip := skipErrors[fname]
		_, isPanic := panicErrors[fname]
		ftype := swigFuncType(fname)

		// Don't alter overloaded versions of panic or skip functions.
		if ftype == overloadFunc {
			split := strings.Split(fname, "__SWIG_")
			oname := split[0]
			if _, ok := skipErrors[oname]; ok {
				goto rename
			} else if _, ok := panicErrors[oname]; ok {
				goto rename
			}
		}

		if !skip && ftype != swigFunc {
			if isPanic && t.Body != nil {
				t.Body.List = append([]ast.Stmt{errReset, errPanic}, t.Body.List...)
			} else if !appendErrResult(t.Type) && t.Body != nil {

				if ftype != swigFunc {
					t.Body.List = append([]ast.Stmt{errStmt}, t.Body.List...)
				}

				n := len(t.Body.List)
				if n > 0 {
					switch rt := t.Body.List[n-1].(type) {
					case *ast.ReturnStmt:
						if len(rt.Results) > 0 {
							rt.Results = append(rt.Results, errRet.(*ast.ReturnStmt).Results...)
						}
					default:
						t.Body.List = append(t.Body.List, &ast.ReturnStmt{})
					}
				}
			}
		}

	rename:
		// Rename functions and methods.
		if v.f != nil {
			if action, ok := v.f[fname]; ok == true {
				if action.Rename == nil {
					break
				}
				r := action.Rename

				if t.Recv != nil && t.Recv.List != nil && len(t.Recv.List) == 1 { // Is it a method?
					switch st := t.Recv.List[0].Type.(type) {
					case *ast.Ident:
						if rename, ok := r.Method[st.Name]; ok == true {
							new := rename.Rename(fname)
							log.Printf("Renaming %v.%v() to %v.%v()", st.Name, fname, st.Name, new)
							t.Name = ast.NewIdent(new)
						}
					}
				} else { // It's a function.
					if r.Function != nil {
						new := r.Function.Rename(fname)
						log.Printf("Renaming %v() to %v()", fname, new)
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
	file, e := parser.ParseFile(fset, fname, nil, 0)
	if e != nil {
		err = e
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
	Function  *Rename
	Method    map[string]*Rename
	Interface map[string]*Rename
}

type Rename struct {
	Prefix   string
	Unexport bool
	Replace  string
}

func (r *Rename) Rename(name string) string {
	if len(r.Replace) > 0 {
		name = r.Replace
	}
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

var errField, anonErrField *ast.Field
var errStmt, errRet, errReset, errPanic ast.Stmt

// initAST initializes the AST for error handling code that will be injected
// into GDAL Go source files generated by SWIG.  For ease of generation of the
// AST, string expressions representing Go source are parsed rather than
// building from scratch.
func initAST(pkg string) (err error) {
	var pstr string
	if len(pkg) > 0 {
		pstr = pkg + "."
	}

	src := fmt.Sprintf(`func (error) (err error) {
defer %sErrorTrap()(&err)
%sErrorReset()
defer func() {
	if err := %sLastError(); err != nil {
		panic(err)
	}
}()
return err;
}`, pstr, pstr, pstr)

	// Generate the AST from the source
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return
	}
	flit := expr.(*ast.FuncLit)

	// Retrieve the named error field for the result.
	errField = flit.Type.Results.List[0]
	// Retrieve the anonymous error field for the result.
	anonErrField = flit.Type.Params.List[0]

	// Retrieve the error trap.
	errStmt = flit.Body.List[0]

	// Retrieve the error reset.
	errReset = flit.Body.List[1]

	// Retrieve the panic statement.
	errPanic = flit.Body.List[2]

	// Retrieve the return statement.
	errRet = flit.Body.List[3]

	return
}

func readErrorFile(name string, store map[string]bool) (err error) {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		store[scanner.Text()] = true
	}

	err = scanner.Err()
	return
}

var skipErrors map[string]bool
var panicErrors map[string]bool

func init() {
	skipErrors = make(map[string]bool)
	panicErrors = make(map[string]bool)
}

var fmods = flag.String("mods", "", "json file defining modifications to functions and methods")
var fmerge = flag.String("merge", "", "go source file defining interfaces to merge")
var fin = flag.String("in", "", "input go source file (defaults to STDIN)")
var fout = flag.String("out", "", "output go source file (defaults to STDOUT)")
var errPkg = flag.String("error-pkg", "cpl", "name of package containing error trap")
var fnoErrors = flag.String("no-errors", "", "file listing function names for which error handling is not required. One function name per line.")
var fpanicErrors = flag.String("panic-errors", "", "file listing function names which should panic on error. One function name per line.")

func main() {
	var err error

	flag.Parse()

	// Initialize the AST used for error injection.
	if err = initAST(*errPkg); err != nil {
		log.Fatal(err)
	}

	// Populate the data structure listing functions for which error handling
	// should be skipped.
	if len(*fnoErrors) > 0 {
		err = readErrorFile(*fnoErrors, skipErrors)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Populate the data structure listing functions for which a panic should be
	// triggered when an error is encountered.
	if len(*fpanicErrors) > 0 {
		err = readErrorFile(*fpanicErrors, panicErrors)
		if err != nil {
			log.Fatal(err)
		}
	}

	var funcs Functions
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
