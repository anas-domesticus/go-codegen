package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParseGoFileOrModule(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gofiles")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gofileContent := `package main
import "fmt"

// MyStruct is a test struct.
// @codegen foo=bar testflag
// {
//     "TestId": 1234
// }
type MyStruct struct {
    Name string ` + "`" + `json:"name"` + "`" + ` // Name of the user
    Age  int    ` + "`" + `json:"age"` + "`" + `  // Age of the user
}

func TestFunc() {
	fmt.Println("Hello World!")
}
`

	err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(gofileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write go file: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
		wantLen  int
	}{
		{
			name:     "Parse single go file",
			filePath: filepath.Join(tempDir, "main.go"),
			wantErr:  false,
			wantLen:  1,
		},
		{
			name:     "Parse directory",
			filePath: tempDir,
			wantErr:  false,
			wantLen:  1,
		},
		{
			name:     "Parse non-existent file",
			filePath: "/non/existent/path",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contexts, err := parseGoFileOrModule(tt.filePath, &Config{
				Name: "test",
			})
			if tt.wantLen != len(contexts) {
				t.Errorf("len(contexts)=%d, wantLen %v", len(contexts), tt.wantLen)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGoFileOrModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
