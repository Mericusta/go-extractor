package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func Test_extractGoMethodMeta(t *testing.T) {
	type args struct {
		extractFilepath string
		structName      string
		methodName      string
	}
	tests := []struct {
		name    string
		args    args
		want    *goMethodMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				extractFilepath: "./testdata/standardProject/pkg/module/module.go",
				structName:      "ExampleStruct",
				methodName:      "ExampleFunc",
			},
			func() *goMethodMeta {
				fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/module/module.go", nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}
				return searchGoMethodMeta(fileAST, "ExampleStruct", "ExampleFunc")
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractGoMethodMeta(tt.args.extractFilepath, tt.args.structName, tt.args.methodName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoMethodMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractGoMethodMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchGoMethodMeta(t *testing.T) {
	type args struct {
		fileAST    *ast.File
		structName string
		methodName string
	}
	tests := []struct {
		name string
		args args
		want *goMethodMeta
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
				methodName: "ExampleFunc",
			},
			func() *goMethodMeta {
				gmm, err := extractGoMethodMeta("./testdata/standardProject/pkg/module/module.go", "ExampleStruct", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gmm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchGoMethodMeta(tt.args.fileAST, tt.args.structName, tt.args.methodName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchGoMethodMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_goMethodMeta_PrintAST(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		gmm  *goMethodMeta
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			"test case 1",
// 			func() *goMethodMeta {
// 				gmm, err := extractGoMethodMeta("./testdata/standardProject/pkg/interface/interface.go", "ExampleInterface", "ExampleFunc")
// 				if err != nil {
// 					panic(err)
// 				}
// 				return gmm
// 			}(),
// 		},
// 		{
// 			"test case 2",
// 			func() *goMethodMeta {
// 				gmm, err := extractGoMethodMeta("./testdata/standardProject/pkg/module/module.go", "ExampleStruct", "ExampleFunc")
// 				if err != nil {
// 					panic(err)
// 				}
// 				return gmm
// 			}(),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.gmm.PrintAST()
// 		})
// 	}
// }
