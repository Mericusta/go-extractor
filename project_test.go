package extractor

import (
	"reflect"
	"testing"
)

func TestExtractGoProjectMeta(t *testing.T) {
	type args struct {
		projectPath string
		ignorePaths map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *GoProjectMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{projectPath: "./testdata/singleFileProject/singleCmd"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
				cmdPath:     ".",
				PackageMap: map[string]*goPackageMeta{
					"main": {
						Name:    "main",
						PkgPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singleCmd",
						pkgFileMap: map[string]*goFileMeta{
							"main.go": nil,
						},
					},
				},
			},
			false,
		},
		{
			"test case 2",
			args{projectPath: "./testdata/singleFileProject/singlePkg"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\singleFileProject\\singlePkg",
				ModuleName:  "pkg",
				pkgPath:     ".",
				PackageMap:  make(map[string]*goPackageMeta),
				ignorePaths: nil,
			},
			false,
		},
		{
			"test case 3",
			args{projectPath: "./testdata/standardProject"},
			&GoProjectMeta{
				ProjectPath: "d:\\Projects\\go-extractor\\testdata\\standardProject",
				ModuleName:  "standard",
				cmdPath:     "cmd",
				pkgPath:     "pkg",
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
