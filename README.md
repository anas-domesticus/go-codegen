# go-codegen

`go-codegen` is a tool designed to parse Go source files and execute Go templates, generating output files based on the data extracted from the source files.

## Installing the tool

To install the tool, run:
```sh
make install
```

## Running the tool

To run the tool, use:
```sh
go-codegen --config go-codegen.yaml
```

## Configuration

The tool is configured using a YAML file with the following structure:

```yaml
---
go_path: src/example.go
template_path: /path/to/template
name: example-template
output_path: generated_foo_%s.go
---
go_path: src/test.go
template_path: /path/to/template
name: foo
output_path: generated_bar_%s.go
remove_fields:
  enable: true
  fields:
    - FieldToRemove
```

Each template is run once per struct found in the source file, with the struct name interpolated into `%s` in the output file name. `go_path` can be a single file or a directory.

A `remove_fields` section will remove any fields with names specified in the array.

## Template Data

The data available to the templates is as follows:

```go
type Field struct {
    SourceName string
    DestName   string
    SourceType string
    DestType   string
    Tags       []map[string]string
    Comment    string
}

type TemplateContext struct {
    TemplateName string
    SourceFile   string
    StructName   string
    Fields       []Field
    Comments     []string
    Config       StructConfig
}

type StructConfig struct {
    Values map[string]string
    Flags  []string
}
```

## Go Module Usage

You can also use `go-codegen` as a Go module by importing the `Templater` type. Inject transformers using the following signature:

```go
type TransformerFn func(*TemplateContext) error
```

Example:

```go
package main

import (
	"fmt"
	codegen "github.com/anas-domesticus/go-codegen/pkg"
	"os"
)

func main() {
	templater := codegen.NewTemplater([]codegen.Config{{
		GoPath:       "src/example/",
		TemplatePath: "example.go.tpl",
		Name:         "example-template",
		OutputPath:   "output_%s.go",
	},
	})
	transformer := func(ctx *codegen.TemplateContext) error {
		ctx.StructName = "OverwrittenStructName"
		return nil
	}
	templater.AddTransformer(transformer)
	err := templater.GenerateFiles()
	if err != nil {
		fmt.Printf("Error running templater: %s\n", err)
		os.Exit(1)
	}
}
```

## Struct Configuration

Structs can be customized using struct tags and comments as shown below:

- including `exclude=true` will result in the struct being skipped
- The struct field tags help customize the output by modifying field names and types in the generated code.

```go
//@codegen exclude=true foo=bar some-flag
type ExampleStruct struct {
    Field1 string `foo:"bar"`
    Field2 int `codegen-new-type:"float64" codegen-new-name:"Field3"`
}
```

## Explanation

1. **Struct-Level Tags:**
    - `// @codegen exclude=true foo=bar some-flag`: These tags placed in a comment above the struct definition provide metadata for the entire struct. They can include:
        - `exclude=true`: Skips the struct during code generation.
        - `foo=bar`: Custom key-value pairs that can be used in templates.
        - `some-flag`: Flags that can be used in templates.

2. **Field-Level Tags:**
    - `foo:"bar"`: Custom tag specific to your application, which can be processed by the code generator.
    - `codegen-new-type:"float64"`: Indicates that the type of this field should be changed to `float64` in the generated code.
    - `codegen-new-name:"Field3"`: Indicates that the field should be renamed to `Field3` in the generated code.

## Example

Given the above struct, it will be presented to the template as:

```go
TemplateContext{
    TemplateName: "example-template",
    SourceFile: "src/example.go",
    StructName: "ExampleStruct",
    Fields: []Field{
        {
            SourceName: "Field1",
            DestName:   "Field1", // defaults to the source name
            SourceType: "string",
            DestType:   "string", // defaults to the source type
            Tags:       []map[string]string{{"foo": "bar"}},
            Comment:    "",
        },
        {
            SourceName: "Field2",
            DestName:   "Field3", // renamed with struct tag
            SourceType: "int",
            DestType:   "float64", // overridden with struct tag
            Tags:       nil,
            Comment:    "",
        },
    },
    Comments: []string{},
    Config: StructConfig{
    Values: map[string]string{"exclude": "true", "foo": "bar"},
    Flags:  []string{"some-flag"},
}
```

## License

MIT License