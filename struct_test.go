package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func Test_extractGoStructMeta(t *testing.T) {
	type args struct {
		extractFilepath string
		structName      string
	}
	tests := []struct {
		name    string
		args    args
		want    *goStructMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				extractFilepath: "./testdata/standardProject/pkg/module/module.go",
				structName:      "ExampleStruct",
			},
			searchGoStructMeta(func() *ast.File {
				fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/module/module.go", nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}
				return fileAST
			}(), "ExampleStruct"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractGoStructMeta(tt.args.extractFilepath, tt.args.structName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoStructMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractGoStructMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchGoStructMeta(t *testing.T) {
	type args struct {
		fileAST    *ast.File
		structName string
	}
	tests := []struct {
		name string
		args args
		want *goStructMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				fileAST: func() *ast.File {
					fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/module/module.go", nil, parser.ParseComments)
					if err != nil {
						panic(err)
					}
					return fileAST
				}(),
				structName: "ExampleStruct",
			},
			func() *goStructMeta {
				gsm, err := extractGoStructMeta("./testdata/standardProject/pkg/module/module.go", "ExampleStruct")
				if err != nil {
					panic(err)
				}
				return gsm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchGoStructMeta(tt.args.fileAST, tt.args.structName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchGoStructMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
