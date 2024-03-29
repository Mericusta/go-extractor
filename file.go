package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func extractGoFilePkgName(fileAbsPath string) (string, error) {
	fileAST, err := parser.ParseFile(token.NewFileSet(), fileAbsPath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return fileAST.Name.String(), nil
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

// CleanFileComment 置空文件中所有注释
func CleanFileComment(r io.Reader) string {
	fileContent, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	isBlock, isComment := false, false
	firstCommentIndex, secondCommentIndex := -1, -1
	builder, commentBuffer := strings.Builder{}, strings.Builder{}
	for index, b := range fileContent {
		switch rune(b) {
		case PunctuationMarkLeftDoubleQuotes:
			if !isComment {
				if !isBlock {
					isBlock = true
				} else {
					isBlock = false
				}
			}
		case '/':
			if !isBlock {
				if firstCommentIndex == -1 {
					firstCommentIndex = index
				} else if secondCommentIndex == -1 {
					secondCommentIndex = index
					isComment = true
					commentBuffer.Reset()
				}
			}
		case '\n':
			if isComment {
				isComment = false
				firstCommentIndex = -1
				secondCommentIndex = -1
				commentBuffer.Reset()
			}
		}

		if !isComment {
			if firstCommentIndex != -1 && secondCommentIndex == -1 {
				if commentBuffer.Len() > 0 {
					// just one /, clear comment buffer
					builder.WriteString(commentBuffer.String())
					builder.WriteByte(b)
					firstCommentIndex = -1
					commentBuffer.Reset()
				} else {
					// first match /
					commentBuffer.WriteByte(b)
				}
			} else {
				builder.WriteByte(b)
			}
		}
	}

	return builder.String()
}
