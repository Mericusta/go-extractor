package extractor

import (
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

func Test_goPackageMeta_SearchStructMeta(t *testing.T) {
	type args struct {
		structName string
	}
	tests := []struct {
		name string
		gpm  *goPackageMeta
		args args
		want *goStructMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goPackageMeta {
				gpm, err := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				if err != nil {
					panic(err)
				}
				return gpm.PackageMap["standardProject/pkg/module"]
			}(),
			args{structName: "ExampleStruct"},
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
			if got := tt.gpm.SearchStructMeta(tt.args.structName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goPackageMeta.SearchStructMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goPackageMeta_SearchInterfaceMeta(t *testing.T) {
	type args struct {
		interfaceName string
	}
	tests := []struct {
		name string
		gpm  *goPackageMeta
		args args
		want *goInterfaceMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goPackageMeta {
				gpm, err := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				if err != nil {
					panic(err)
				}
				return gpm.PackageMap["standardProject/pkg/pkgInterface"]
			}(),
			args{interfaceName: "ExampleInterface"},
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
			if got := tt.gpm.SearchInterfaceMeta(tt.args.interfaceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goPackageMeta.SearchInterfaceMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goPackageMeta_SearchFunctionMeta(t *testing.T) {
	type args struct {
		functionName string
	}
	tests := []struct {
		name string
		gpm  *goPackageMeta
		args args
		want *goFunctionMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goPackageMeta {
				gpm, err := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				if err != nil {
					panic(err)
				}
				return gpm.PackageMap["standardProject/pkg"]
			}(),
			args{functionName: "ExampleFunc"},
			func() *goFunctionMeta {
				gsm, err := extractGoFunctionMeta("./testdata/standardProject/pkg/pkg.go", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gsm
			}(),
		},
		{
			"test case 2",
			func() *goPackageMeta {
				gpm, err := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				if err != nil {
					panic(err)
				}
				return gpm.PackageMap["standardProject/pkg/module"]
			}(),
			args{functionName: "ExampleFunc"},
			func() *goFunctionMeta {
				gsm, err := extractGoFunctionMeta("./testdata/standardProject/pkg/module/module.go", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gsm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gpm.SearchFunctionMeta(tt.args.functionName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goPackageMeta.SearchFunctionMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goPackageMeta_SearchMethodMeta(t *testing.T) {
	type args struct {
		structName string
		methodName string
	}
	tests := []struct {
		name string
		gpm  *goPackageMeta
		args args
		want *goMethodMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			func() *goPackageMeta {
				gpm, err := ExtractGoProjectMeta("./testdata/standardProject", map[string]struct{}{
					"./testdata/standardProject/vendor": {},
				})
				if err != nil {
					panic(err)
				}
				return gpm.PackageMap["standardProject/pkg/module"]
			}(),
			args{
				structName: "ExampleStruct",
				methodName: "ExampleFunc",
			},
			func() *goMethodMeta {
				gsm, err := extractGoMethodMeta("./testdata/standardProject/pkg/module/module.go", "ExampleStruct", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gsm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gpm.SearchMethodMeta(tt.args.structName, tt.args.methodName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goPackageMeta.SearchMethodMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
