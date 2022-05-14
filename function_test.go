package extractor

import "testing"

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
