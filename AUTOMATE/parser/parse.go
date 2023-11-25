package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "/Users/bfuzzy1/Code/tests/tests" // Set your directory path here
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			yamlOutput, err := parseGoFile(path)
			if err != nil {
				return err
			}
			fmt.Println(yamlOutput) // Print YAML output for each file
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}
}

func parseGoFile(path string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	return generateYAMLRepresentation(node), nil
}

func generateYAMLRepresentation(node *ast.File) string {
	var yamlBuilder strings.Builder
	for _, d := range node.Decls {
		if fn, isFn := d.(*ast.FuncDecl); isFn {
			yamlBuilder.WriteString("- function: " + fn.Name.Name + "\n")
			yamlBuilder.WriteString("  body:\n")
			generateFunctionBodyYAML(fn.Body, &yamlBuilder, "    ")
		}
	}
	return yamlBuilder.String()
}

func generateFunctionBodyYAML(block *ast.BlockStmt, yamlBuilder *strings.Builder, indent string) {
	if block == nil {
		return
	}
	for _, stmt := range block.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			yamlBuilder.WriteString(indent + "if:\n")
			generateFunctionBodyYAML(s.Body, yamlBuilder, indent+"  ")
			if s.Else != nil {
				// Check the type of the Else statement
				switch elseStmt := s.Else.(type) {
				case *ast.BlockStmt:
					yamlBuilder.WriteString(indent + "else:\n")
					generateFunctionBodyYAML(elseStmt, yamlBuilder, indent+"  ")
				case *ast.IfStmt:
					// Handle else if
					yamlBuilder.WriteString(indent + "else if:\n")
					generateFunctionBodyYAML(elseStmt.Body, yamlBuilder, indent+"  ")
				}
			}
			// Add cases for other statement types as needed
		}
	}
}
