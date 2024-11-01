package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func parseGoFileOrModule(path string, cfg *Config) ([]TemplateContext, error) {
	var files []string

	// Collect Go files either from a single file or from a directory
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		// Walk the directory to find all Go files
		err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Single Go file
		files = append(files, path)
	}

	var contexts []TemplateContext
	fset := token.NewFileSet()

	for _, filePath := range files {
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {

				// Parsing comments first as we may have an exclude flag
				var comments []string
				var structConfig StructConfig
				if genDecl.Doc != nil {
					for _, comment := range genDecl.Doc.List {
						if strings.Contains(comment.Text, "@codegen") {
							structConfig = parseCommentAndLoadConfig(comment.Text)
						}
						comments = append(comments, comment.Text)
					}
					v, ok := structConfig.Values["exclude"]
					if ok && v == "true" {
						continue
					}
				}

				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				var fields []Field
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fieldType := ExtractFieldName(field.Type)
						fieldTag := ""
						if field.Tag != nil {
							fieldTag = field.Tag.Value
						}

						fieldComment := ""
						if field.Doc != nil {
							fieldComment = field.Doc.Text()
						}
						newField := Field{
							SourceName: Name(name.Name),
							DestName:   Name(name.Name),
							Counter:    len(fields),
							SourceType: fieldType,
							DestType:   fieldType,
							Tags:       parseTags(fieldTag),
							Comment:    fieldComment,
						}
						// Checking for exclusion tag on field
						for i := range newField.Tags {
							v, ok := newField.Tags[i]["codegen-exclude"]
							if ok && v == "true" {
								continue
							}
						}
						fields = append(fields, newField)
					}
				}

				contexts = append(contexts, TemplateContext{
					TemplateName: cfg.Name,
					SourceFile:   filePath,
					StructName:   Name(typeSpec.Name.Name),
					Fields:       fields,
					Comments:     comments,
					Config:       structConfig,
				})
			}
		}
	}

	return contexts, nil
}

// ExtractFieldName extracts the field name from an ast.Expr using reflection
func ExtractFieldName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		pkgName := ExtractFieldName(v.X)
		return pkgName + "." + ExtractFieldName(v.Sel)
	case *ast.StarExpr:
		return "*" + ExtractFieldName(v.X)
	case *ast.ArrayType:
		return "[]" + ExtractFieldName(v.Elt)
	case *ast.MapType:
		return "map[" + ExtractFieldName(v.Key) + "]" + ExtractFieldName(v.Value)
	// Add other cases for different ast.Expr types as needed
	default:
		rv := reflect.ValueOf(expr)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			if nameField := rv.FieldByName("Name"); nameField.IsValid() && nameField.Kind() == reflect.String {
				return nameField.String()
			}
		}
	}
	return ""
}
