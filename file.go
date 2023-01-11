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
	fileSet *token.FileSet
	fileAST *ast.File
	Name    string // filename
	Path    string // file path
}

func ExtractGoFileMeta(extractFilepath string) (*GoFileMeta, error) {
	fileSet := token.NewFileSet()
	fileAST, err := parser.ParseFile(fileSet, extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, obj := range fileAST.Scope.Objects {
		switch obj.Kind {
		}
	}

	fileMeta := &GoFileMeta{
		fileSet: fileSet,
		fileAST: fileAST,
		Name:    filepath.Base(extractFilepath),
		Path:    extractFilepath,
	}

	return fileMeta, nil
}

func (gfm *GoFileMeta) PrintAST() {
	ast.Print(gfm.fileSet, gfm.fileAST)
}

func (gfm *GoFileMeta) OutputAST() {
	outputFile, err := os.OpenFile(fmt.Sprintf("%v.ast", gfm.Path), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	ast.Fprint(outputFile, gfm.fileSet, gfm.fileAST, ast.NotNilFilter)
}

func (gfm *GoFileMeta) PkgName() string {
	return gfm.fileAST.Name.Name
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
