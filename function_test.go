package extractor

import (
	"runtime"
	"testing"
)

func TestGoFunctionDeclaration_Traversal(t *testing.T) {
	type args struct {
		deep int
	}
	tests := []struct {
		name string
		d    *GoFunctionDeclaration
		args args
	}{
		// TODO: Add test cases.
		{
			name: "UnitTestBundle.go ReadUnitTestFile Declaration test",
			d: func() *GoFunctionDeclaration {
				fdMap := ExtractGoFileFunctionDeclaration(ReadUnitTestFile("UnitTestBundle.go"))
				return fdMap["ReadUnitTestFile"]
			}(),
			args: args{
				deep: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Traversal(tt.args.deep)
		})
	}
}

func TestGoFunctionDeclaration_MakeUp(t *testing.T) {
	tests := []struct {
		name string
		d    *GoFunctionDeclaration
		want string
	}{
		// TODO: Add test cases.
		{
			name: "UnitTestBundle.go ReadUnitTestFile Declaration test",
			d: func() *GoFunctionDeclaration {
				fdMap := ExtractGoFileFunctionDeclaration(ReadUnitTestFile("UnitTestBundle.go"))
				return fdMap["ReadUnitTestFile"]
			}(),
			want: func() string {
				if runtime.GOOS == "windows" {
					return `func ReadUnitTestFile(p string) []byte {` + "\r" + `
	c, err := os.ReadFile(p)` + "\r" + `
	if err != nil {` + "\r" + `
		panic(err)` + "\r" + `
	}` + "\r" + `
	return c` + "\r" + `
}`
				}
				return `func ReadUnitTestFile(p string) []byte {
	c, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return c
}`
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MakeUp(); got != tt.want {
				t.Errorf("GoFunctionDeclaration.MakeUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
