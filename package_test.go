package extractor

import "testing"

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
