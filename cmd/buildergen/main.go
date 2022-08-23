package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

const (
	fileNameFlag = "f"
	goExt        = ".go"
)

var originalFileName string

func init() {
	flag.StringVar(&originalFileName, fileNameFlag, "", "go file for generating builders")
	flag.Parse()
}

func main() {
	if originalFileName == "" {
		logger.Fatalf("filename is empty. Use -%s for set", fileNameFlag)
	}
	if len(originalFileName) < 4 || originalFileName[len(originalFileName)-3:] != goExt {
		logger.Fatalf("file isn't go format. Use go file")
	}

	f, err := parser.ParseFile(token.NewFileSet(), originalFileName, nil, parser.AllErrors)
	if err != nil {
		logger.Fatalf("parse error: %s", err.Error())
	}

	p := Package{}

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			p.Name = n.Name.Name
		case *ast.TypeSpec:
			s := Struct{}
			s.Name = n.Name.Name

			switch n.Type.(type) {
			case *ast.StructType:
				structType := n.Type.(*ast.StructType)

				// Iterating fields
				for _, field := range structType.Fields.List {
					if len(field.Names) == 1 {
						var fieldTypeName string
						var isPointer bool

						switch fieldType := field.Type.(type) {
						// Value
						case *ast.Ident:
							fieldTypeName = fieldType.Name
							isPointer = false
						// Pointer
						case *ast.StarExpr:
							fieldTypeName = fieldType.X.(*ast.Ident).Name
							isPointer = true
						}

						// TODO: Maybe there is a better option to get privacy :)
						rs := []rune(field.Names[0].Name)
						if len(rs) == 0 || rs[0] < 'A' || rs[0] > 'Z' {
							continue
						}

						if isPointer {
							fieldTypeName = "*" + fieldTypeName
						}

						s.Fields = append(s.Fields, Field{
							Name: field.Names[0].Name,
							Type: fieldTypeName,
						})
					}
				}
			}

			p.Structs = append(p.Structs, s)
		}

		return true
	})

	t, err := template.New("builder").Parse(templ)
	if err != nil {
		logger.Fatalf("templater error: %s", err.Error())
	}

	b := &bytes.Buffer{}
	if err = t.Execute(b, p); err != nil {
		logger.Fatalf("execute template error: %s", err.Error())
	}

	fn := strings.TrimRight(filepath.Base(originalFileName), goExt) + "_builder" + goExt

	ioutil.WriteFile(filepath.Join(filepath.Dir(originalFileName), fn), b.Bytes(), 0o664)
}
