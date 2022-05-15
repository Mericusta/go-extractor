package extractor

import (
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
			name: "UnitTestBundle.go GoFunctionDeclaration.Traversal test",
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
			name: "UnitTestBundle.go GoFunctionDeclaration.MakeUp test",
			d: func() *GoFunctionDeclaration {
				fdMap := ExtractGoFileFunctionDeclaration(ReadUnitTestFile("UnitTestBundle.go"))
				return fdMap["ReadUnitTestFile"]
			}(),
			want: `
func ReadUnitTestFile(p string) []byte {
	c, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return c
}
`,
		},
		{
			name: "file.go GoFmtFile Declaration test",
			d: func() *GoFunctionDeclaration {
				fdMap := ExtractGoFileFunctionDeclaration(ReadUnitTestFile("file.go"))
				return fdMap["GoFmtFile"]
			}(),
			want: `
func GoFmtFile(p string) {
	if _, err := os.Stat(p); !(err == nil || os.IsExist(err)) {
		panic(fmt.Sprintf("%v not exist", p))
	}
	cmd := exec.Command("go", "fmt", p)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		panic(cmd.Stderr.(*bytes.Buffer).String())
	}
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MakeUp(); got != tt.want {
				t.Errorf("GoFunctionDeclaration.MakeUp() = %v, want %v", got, tt.want)
				t.Logf("got: %v", []byte(got))
				t.Logf("tt.want: %v", []byte(tt.want))
			}
		})
	}
}
