package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func Test_extractGoInterfaceMeta(t *testing.T) {
	type args struct {
		extractFilepath string
		interfaceName   string
	}
	tests := []struct {
		name    string
		args    args
		want    *goInterfaceMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				extractFilepath: "./testdata/standardProject/pkg/interface/interface.go",
				interfaceName:   "ExampleInterface",
			},
			searchGoInterfaceMeta(func() *ast.File {
				fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/interface/interface.go", nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}
				return fileAST
			}(), "ExampleInterface"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractGoInterfaceMeta(tt.args.extractFilepath, tt.args.interfaceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoInterfaceMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractGoInterfaceMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchGoInterfaceMeta(t *testing.T) {
	type args struct {
		fileAST       *ast.File
		interfaceName string
	}
	tests := []struct {
		name string
		args args
		want *goInterfaceMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				fileAST: func() *ast.File {
					fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/interface/interface.go", nil, parser.ParseComments)
					if err != nil {
						panic(err)
					}
					return fileAST
				}(),
				interfaceName: "ExampleInterface",
			},
			func() *goInterfaceMeta {
				gim, err := extractGoInterfaceMeta("./testdata/standardProject/pkg/interface/interface.go", "ExampleInterface")
				if err != nil {
					panic(err)
				}
				return gim
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchGoInterfaceMeta(tt.args.fileAST, tt.args.interfaceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchGoInterfaceMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goInterfaceMeta_SearchMethodDecl(t *testing.T) {
	type args struct {
		methodName string
	}
	tests := []struct {
		name string
		gim  *goInterfaceMeta
		args args
		want *goMethodMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goInterfaceMeta {
				gim, err := extractGoInterfaceMeta("./testdata/standardProject/pkg/interface/interface.go", "ExampleInterface")
				if err != nil {
					panic(err)
				}
				return gim
			}(),
			args{methodName: "ExampleFunc"},
			func() *goMethodMeta {
				gmm, err := extractGoMethodMeta("./testdata/standardProject/pkg/interface/interface.go", "ExampleInterface", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gmm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gim.SearchMethodDecl(tt.args.methodName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goInterfaceMeta.SearchMethodDecl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goInterfaceMeta_ForeachMethodDecl(t *testing.T) {
	type args struct {
		f func(*ast.Field) bool
	}
	tests := []struct {
		name string
		gim  *goInterfaceMeta
		args args
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goInterfaceMeta {
				gim, err := extractGoInterfaceMeta("./testdata/standardProject/pkg/interface/interface.go", "ExampleInterface")
				if err != nil {
					panic(err)
				}
				return gim
			}(),
			args{func(f *ast.Field) bool {
				if f.Doc != nil {
					fmt.Printf("%v\n", f.Doc.Text())
				}
				return true
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gim.ForeachMethodDecl(tt.args.f)
		})
	}
}
