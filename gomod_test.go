package extractor

import "testing"

func Test_extractGoModuleName(t *testing.T) {
	type args struct {
		goModFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{goModFilePath: "go.mod"},
			"github.com/Mericusta/go-extractor",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractGoModuleName(tt.args.goModFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoModuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractGoModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
