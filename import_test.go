package extractor

import (
	"testing"
)

func TestGoImportDeclaration_Traversal(t *testing.T) {
	type args struct {
		deep int
	}
	tests := []struct {
		name string
		d    *GoImportDeclaration
		args args
	}{
		// TODO: Add test cases.
		{
			name: "UnitTestBundle.go GoImportDeclaration.Traversal test",
			d:    ExtractGoFileImportDeclaration(ReadUnitTestFile("UnitTestBundle.go")),
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

func TestGoImportDeclaration_MakeUp(t *testing.T) {
	type args struct {
		withBracket bool
	}
	tests := []struct {
		name string
		d    *GoImportDeclaration
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "UnitTestBundle.go ImportDeclaration.MakeUp test",
			d:    ExtractGoFileImportDeclaration(ReadUnitTestFile("UnitTestBundle.go")),
			args: args{
				withBracket: true,
			},
			want: `
import (
	windowsOS "os"
)
`,
		},
		{
			name: "UnitTestBundle.go ImportDeclaration.MakeUp test",
			d:    ExtractGoFileImportDeclaration(ReadUnitTestFile("UnitTestBundle.go")),
			args: args{
				withBracket: false,
			},
			want: `
import windowsOS "os"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MakeUp(tt.args.withBracket); got != tt.want {
				t.Errorf("GoImportDeclaration.MakeUp() = %v, want %v", got, tt.want)
				t.Logf("got %v\n", []byte(got))
				t.Logf("want %v\n", []byte(tt.want))
			}
		})
	}
}
