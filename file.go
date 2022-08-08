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

type GoFileInfo struct {
	fileSet                 *token.FileSet
	fileAST                 *ast.File
	Name                    string                                      // 文件名
	Path                    string                                      // 相对项目根目录的路径
	ImportStruct            map[string]map[string]struct{}              // 该文件引入的外部包
	StructDefinitionMap     map[string]map[string]*GoVariableDefinition // 该文件定义的结构体
	InterfaceDeclarationMap map[string]*GoInterfaceInfo                 // 该文件定义的接口
}

type goFileMeta struct {
	fileSet *token.FileSet
	fileAST *ast.File
	Name    string // filename
	Path    string // file path
}

func extractGoFileMeta(extractFilepath string) (*goFileMeta, error) {
	fileSet := token.NewFileSet()
	fileAST, err := parser.ParseFile(fileSet, extractFilepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, obj := range fileAST.Scope.Objects {
		switch obj.Kind {
		}
	}

	fileMeta := &goFileMeta{
		fileSet: fileSet,
		fileAST: fileAST,
		Name:    filepath.Base(extractFilepath),
		Path:    extractFilepath,
	}

	return fileMeta, nil
}

func (gfm *goFileMeta) PkgName() string {
	return gfm.fileAST.Name.Name
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

// go fmt 格式化文件
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
