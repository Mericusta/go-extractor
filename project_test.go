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
		{
			"test case 1: single main package file",
			args{projectPath: "./testdata/singleFileProject/singleCmd/main.go"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd\\main.go",
				PackageMap: map[string]*goPackageMeta{
					"main": {
						Name:    "main",
						PkgPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
						pkgFileMap: map[string]*goFileMeta{
							"main.go": func() *goFileMeta {
								gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd\\main.go")
								return gfm
							}(),
						},
					},
				},
				ignorePaths: make(map[string]struct{}),
			},
			false,
		},
		{
			"test case 2: single pkg file",
			args{projectPath: "./testdata/singleFileProject/singlePkg/pkg.go"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg\\pkg.go",
				PackageMap: map[string]*goPackageMeta{
					"pkg": {
						Name:    "pkg",
						PkgPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
						pkgFileMap: map[string]*goFileMeta{
							"pkg.go": func() *goFileMeta {
								gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg\\pkg.go")
								return gfm
							}(),
						},
					},
				},
				ignorePaths: make(map[string]struct{}),
			},
			false,
		},
		{
			"test case 3",
			args{projectPath: "./testdata/singleFileProject/singleCmd"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
				ModuleName:  "singleCmd",
				PackageMap: map[string]*goPackageMeta{
					"main": {
						Name:    "main",
						PkgPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
						pkgFileMap: map[string]*goFileMeta{
							"main.go": func() *goFileMeta {
								gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd\\main.go")
								return gfm
							}(),
						},
					},
				},
				ignorePaths: make(map[string]struct{}),
			},
			false,
		},
		{
			"test case 4",
			args{projectPath: "./testdata/singleFileProject/singlePkg"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
				ModuleName:  "singlePkg",
				PackageMap: map[string]*goPackageMeta{
					"singlePkg/pkg": {
						Name:       "pkg",
						PkgPath:    "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
						ImportPath: "singlePkg/pkg",
						pkgFileMap: map[string]*goFileMeta{
							"pkg.go": func() *goFileMeta {
								gfm, _ := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg\\pkg.go")
								return gfm
							}(),
						},
					},
				},
				ignorePaths: make(map[string]struct{}),
			},
			false,
		},
		{
			"test case 5",
			args{
				projectPath: "./testdata/standardProject",
				ignorePaths: map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				},
			},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\standardProject",
				ModuleName:  "standardProject",
				PackageMap: map[string]*goPackageMeta{
					"main": {
						Name:    "main",
						PkgPath: "d:\\Projects\\go-extractor\\testdata\\standardProject\\cmd",
						pkgFileMap: map[string]*goFileMeta{
							"init.go": func() *goFileMeta {
								gfm, err := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\cmd\\init.go")
								if err != nil {
									panic(err)
								}
								return gfm
							}(),
							"main.go": func() *goFileMeta {
								gfm, err := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\cmd\\main.go")
								if err != nil {
									panic(err)
								}
								return gfm
							}(),
						},
					},
					"standardProject/pkg": {
						Name:       "pkg",
						PkgPath:    "d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg",
						ImportPath: "standardProject/pkg",
						pkgFileMap: map[string]*goFileMeta{
							"pkg.go": func() *goFileMeta {
								gfm, err := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg\\pkg.go")
								if err != nil {
									panic(err)
								}
								return gfm
							}(),
						},
					},
					"standardProject/pkg/module": {
						Name:       "module",
						PkgPath:    "d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg\\module",
						ImportPath: "standardProject/pkg/module",
						pkgFileMap: map[string]*goFileMeta{
							"module.go": func() *goFileMeta {
								gfm, err := extractGoFileMeta("d:\\Projects\\go-extractor\\testdata\\standardProject\\pkg\\module\\module.go")
								if err != nil {
									panic(err)
								}
								return gfm
							}(),
						},
					},
				},
				ignorePaths: map[string]struct{}{
					"d:\\Projects\\go-extractor\\testdata\\standardProject\\vendor": {},
				},
			},
			false,
		},
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
