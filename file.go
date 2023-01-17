package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
)

type GoFileMeta struct {
	*meta
	fileSet *token.FileSet
}

func ExtractGoFileMeta(extractFilepath string) (*GoFileMeta, error) {
	fileAbsPath, err := filepath.Abs(extractFilepath)
	if err != nil {
		return nil, err
	}

	fileSet := token.NewFileSet()
	fileAST, err := parser.ParseFile(fileSet, fileAbsPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	meta := &GoFileMeta{
		meta: &meta{
			node: fileAST,
			name: filepath.Base(fileAbsPath),
			path: fileAbsPath,
		},
		fileSet: fileSet,
	}

	return meta, nil
}

func (gfm *GoFileMeta) OutputAST() {
	outputFile, err := os.OpenFile(fmt.Sprintf("%v.ast", gfm.path), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	ast.Fprint(outputFile, gfm.fileSet, gfm.node, ast.NotNilFilter)
}

func (gfm *GoFileMeta) Name() string {
	return gfm.name
}

func (gfm *GoFileMeta) Path() string {
	return gfm.path
}

func (gfm *GoFileMeta) PkgName() string {
	return gfm.node.(*ast.File).Name.String()
}

// GoFmtFile go fmt 格式化文件
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
