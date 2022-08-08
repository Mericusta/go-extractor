package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
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
			}
		})
	}
}

func TestExtractGoFileFunctionDeclaration(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]*GoFunctionDeclaration
	}{
		// TODO: Add test cases.
		{
			"example test",
			args{
				content: []byte(`
func TestNewPoint(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		// TODO: Add test cases.
		{
			"Point ❤",
			args{r: '❤'},
			Point{
				ShapeContext: ShapeContext{
					BasicContext: core.NewBasicContext(core.Size{
						Width: 2, Height: 1,
					}),
				},
				r: '❤',
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPoint(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
`),
			},
			map[string]*GoFunctionDeclaration{
				"TestNewPoint": {
					FunctionSignature: "TestNewPoint",
					This:              nil,
					ParamsList:        ExtractorFunctionParamsList([]byte("t *testing.T")),
					ReturnList:        nil,
					BodyContent: []byte(`
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		// TODO: Add test cases.
		{
			"Point ❤",
			args{r: '❤'},
			Point{
				ShapeContext: ShapeContext{
					BasicContext: core.NewBasicContext(core.Size{
						Width: 2, Height: 1,
					}),
				},
				r: '❤',
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPoint(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
`),
					Content: []byte(`
func TestNewPoint(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		// TODO: Add test cases.
		{
			"Point ❤",
			args{r: '❤'},
			Point{
				ShapeContext: ShapeContext{
					BasicContext: core.NewBasicContext(core.Size{
						Width: 2, Height: 1,
					}),
				},
				r: '❤',
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPoint(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractGoFileFunctionDeclaration(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractGoFileFunctionDeclaration() = %v, want %v", got, tt.want)
				t.Logf("got |%+v|", string(got["TestNewPoint"].Content))
				t.Logf("want |%+v|", string(tt.want["TestNewPoint"].Content))
			}
		})
	}
}

func Test_extractGoFunctionMeta(t *testing.T) {
	type args struct {
		extractFilepath string
		functionName    string
	}
	tests := []struct {
		name    string
		args    args
		want    *goFunctionMeta
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				extractFilepath: "./testdata/standardProject/pkg/pkg.go",
				functionName:    "ExampleFunc",
			},
			searchGoFunctionMeta(func() *ast.File {
				fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/pkg.go", nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}
				return fileAST
			}(), "ExampleFunc"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractGoFunctionMeta(tt.args.extractFilepath, tt.args.functionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractGoFunctionMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractGoFunctionMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchGoFunctionMeta(t *testing.T) {
	type args struct {
		fileAST      *ast.File
		functionName string
	}
	tests := []struct {
		name string
		args args
		want *goFunctionMeta
	}{
		// TODO: Add test cases.
		{
			"test case 1",
			args{
				fileAST: func() *ast.File {
					fileAST, err := parser.ParseFile(token.NewFileSet(), "./testdata/standardProject/pkg/pkg.go", nil, parser.ParseComments)
					if err != nil {
						panic(err)
					}
					return fileAST
				}(),
				functionName: "ExampleFunc",
			},
			func() *goFunctionMeta {
				gfm, err := extractGoFunctionMeta("./testdata/standardProject/pkg/pkg.go", "ExampleFunc")
				if err != nil {
					panic(err)
				}
				return gfm
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchGoFunctionMeta(tt.args.fileAST, tt.args.functionName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchGoFunctionMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}
