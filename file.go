package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GoFileMetaTypeConstraints interface {
	*ast.File

	ast.Node
}

// GoFileMeta go 文件 的 meta 数据
type GoFileMeta[T GoFileMetaTypeConstraints] struct {
	// 组合基本 meta 数据
	// ast 节点，要求为 *ast.File
	// 以 ast 节点 为单位执行 AST/PrintAST/Expression/Format
	*meta[T]

	// 文件 token set 集合
	fileSet *token.FileSet

	// 文件名称
	ident string

	// 文件包名称
	packageName string
}

// NewGoFileMeta 构造 go 文件 的 meta 数据
func NewGoFileMeta[T GoFileMetaTypeConstraints](m *meta[T], fs *token.FileSet, fn string) *GoFileMeta[T] {
	return &GoFileMeta[T]{meta: m, fileSet: fs, ident: fn}
}

// -------------------------------- extractor --------------------------------

// ExtractGoFileMeta 通过文件的绝对路径提取文件的 meta 数据
func ExtractGoFileMeta[T GoFileMetaTypeConstraints](extractFilepath string) (*GoFileMeta[T], error) {
	fileAbsPath, err := filepath.Abs(extractFilepath)
	if err != nil {
		return nil, err
	}

	fileSet := token.NewFileSet()
	fileAST, err := parser.ParseFile(fileSet, fileAbsPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	meta := &GoFileMeta[T]{
		// meta:        &meta{node: fileAST, path: fileAbsPath},
		meta:        newMeta[T](fileAST, fileAbsPath),
		fileSet:     fileSet,
		ident:       filepath.Base(fileAbsPath),
		packageName: fileAST.Name.String(),
	}

	return meta, nil
}

// -------------------------------- extractor --------------------------------

// OutputAST 在文件所属的目录下创建一个 同名+.ast 后缀的文件，输出该文件的 ast 树
func (gfm *GoFileMeta[T]) OutputAST() {
	outputFile, err := os.OpenFile(fmt.Sprintf("%v.ast", gfm.path), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	ast.Fprint(outputFile, gfm.fileSet, gfm.node, ast.NotNilFilter)
}

// -------------------------------- unit test --------------------------------

func (gfm *GoFileMeta[T]) Ident() string       { return gfm.ident }
func (gfm *GoFileMeta[T]) PackageName() string { return gfm.packageName }

// -------------------------------- unit test --------------------------------

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
	fileContent, err := io.ReadAll(r)
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
