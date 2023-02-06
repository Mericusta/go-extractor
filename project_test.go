package extractor

import (
	"reflect"
	"testing"
)

func TestExtractGoProjectMetaByDir(t *testing.T) {
	type args struct {
		projectPath   string
		toHandlePaths map[string]struct{}
		spec          bool
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
			args{
				projectPath:   standardProjectAbsPath,
				toHandlePaths: nil,
				spec:          false,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractGoProjectMetaByDir(tt.args.projectPath, tt.args.toHandlePaths, tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractGoProjectMetaByDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractGoProjectMetaByDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
