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
```

Each template is run once per struct found in the source file, with the struct name interpolated into `%s` in the output file name. `go_path` can be a single file or a directory.

## Template Data

The data available to the templates is as follows:

```go
type Field struct {
    Name    string
    Type    string
    Tags    []map[string]string
    Comment string
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

You can also use `go-codegen` as a Go module by importing the `Templater` type. Inject transformers using the following
signature:

```go
type TransformerFn func(*TemplateContext) error
```

## Struct Configuration

Structs can be customized using struct tags and comments as shown below:

- including `exclude=true` will result in the struct being skipped

```go
//@codegen exclude=true foo=bar some-flag
type ExampleStruct struct {
    Field1 string `foo:"bar"`
    Field2 int
}
```

## Example

Given the above struct, it will be presented to the template as:

```go
TemplateContext{
    TemplateName: "example-template",
    SourceFile: "src/example.go",
    StructName: "ExampleStruct",
    Fields: []Field{
        {
            Name:    "Field1",
            Type:    "string",
            Tags:    []map[string]string{{"foo": "bar"}},
            Comment: "",
        },
        {
            Name:    "Field2",
            Type:    "int",
            Tags:    nil,
            Comment: "",
        },
    },
        Comments: []string{},
        Config: StructConfig{
        Values: map[string]string{"exclude": "true", "foo": "bar"},
        Flags:  []string{"some-flag"},
    },
}
```

## License

MIT License