package extractor

import (
	"reflect"
	"testing"
)

func TestExtractGoFileMeta(t *testing.T) {
	type args struct {
		extractFilepath string
	}
	tests := []struct {
		name    string
		args    args
		want    *GoFileMeta
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractGoFileMeta(tt.args.extractFilepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractGoFileMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractGoFileMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
