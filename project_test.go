package extractor

import (
	"reflect"
	"testing"
)

func TestExtractGoProjectMeta(t *testing.T) {
	type args struct {
		projectPath string
		ignorePaths map[string]struct{}
	}
	tests := []struct {
		name    string
		args    args
		want    *GoProjectMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		// {
		// 	"test case 1",
		// 	args{projectPath: "./testdata/singleFileProject/singleCmd"},
		// 	&GoProjectMeta{
		// 		ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
		// 		ModuleName:  "singleCmd",
		// 		PackageMap: map[string]*goPackageMeta{
		// 			"main": {
		// 				Name:    "main",
		// 				PkgPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
		// 				pkgFileMap: map[string]*goFileMeta{
		// 					"main.go": func() *goFileMeta {
		// 						gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd\\main.go")
		// 						return gfm
		// 					}(),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	false,
		// },
		// {
		// 	"test case 2",
		// 	args{projectPath: "./testdata/singleFileProject/singlePkg"},
		// 	&GoProjectMeta{
		// 		ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
		// 		ModuleName:  "singlePkg",
		// 		PackageMap: map[string]*goPackageMeta{
		// 			"singlePkg/pkg": {
		// 				Name:       "pkg",
		// 				PkgPath:    "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
		// 				ImportPath: "singlePkg/pkg",
		// 				pkgFileMap: map[string]*goFileMeta{
		// 					"pkg.go": func() *goFileMeta {
		// 						gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg\\pkg.go")
		// 						return gfm
		// 					}(),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	false,
		// },
		// {
		// 	"test case 3",
		// 	args{
		// 		projectPath: "./testdata/standardProject",
		// 		ignorePaths: map[string]struct{}{
		// 			"./testdata/standardProject/vendor": {},
		// 		},
		// 	},
		// 	&GoProjectMeta{
		// 		ProjectPath: "d:\\Projects\\go-extractor\\testdata\\standardProject",
		// 		ModuleName:  "standard",
		// 		PackageMap: map[string]*goPackageMeta{
		// 			"main": {
		// 				Name:    "main",
		// 				PkgPath: "d:\\Projects\\go-extractor\\testdata\\standardProject\\cmd",
		// 				pkgFileMap: map[string]*goFileMeta{
		// 					"init.go": func() *goFileMeta {
		// 						gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\init.go")
		// 						return gfm
		// 					}(),
		// 					"main.go": func() *goFileMeta {
		// 						gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\main.go")
		// 						return gfm
		// 					}(),
		// 				},
		// 			},
		// 			"standard/pkg": {
		// 				Name:       "pkg",
		// 				PkgPath:    "d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg",
		// 				ImportPath: "standard/pkg",
		// 				pkgFileMap: map[string]*goFileMeta{
		// 					"pkg.go": func() *goFileMeta {
		// 						gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg\\pkg.go")
		// 						return gfm
		// 					}(),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractGoProjectMeta(tt.args.projectPath, tt.args.ignorePaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractGoProjectMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractGoProjectMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
