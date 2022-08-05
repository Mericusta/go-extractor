package extractor

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestGoPackageInfo_MakeUp(t *testing.T) {
	tests := []struct {
		name string
		gpi  *GoPackageInfo
		want string
	}{
		// TODO: Add test cases.
		{
			name: "UnitTestBundle.go Package Info test",
			gpi: &GoPackageInfo{
				Name: ExtractGoFilePackage(ReadUnitTestFile("UnitTestBundle.go")),
			},
			want: "package extractor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gpi.MakeUp(); got != tt.want {
				t.Errorf("GoPackageInfo.MakeUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goPackageMeta_SearchStructDeclaration(t *testing.T) {
	tests := []struct {
		name string
		gpm  *goPackageMeta
		want *ast.StructType
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goPackageMeta {
				gpm, _ := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				return gpm.PackageMap["standardProject/pkg"]
			}(),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gpm.SearchStructDeclaration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goPackageMeta.SearchStructDeclaration() = %v, want %v", got, tt.want)
			}
		})
	}
}
