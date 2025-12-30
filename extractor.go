package extractor

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strings"
)

// Meta 所有 meta 的接口
type Meta interface {
	AST() []byte
	PrintAST()
	Expression() string
}

// meta 基本 meta 数据
type meta struct {
	// 当前 meta 的 ast 节点
	node ast.Node

	// 当前 meta 的 ast 节点所属的文件的绝对路径
	path string
}

func newMeta(node ast.Node, path string) *meta {
	return &meta{node: node, path: path}
}

// copyMeta 保持 path 不变的情况下构造新 ast 节点的 meta 数据
func (m *meta) copyMeta(node ast.Node) *meta {
	return &meta{node: node, path: m.path}
}

// AST 获取当前 meta 的 ast 节点树
func (m *meta) AST() []byte {
	jsonAst, err := json.MarshalIndent(m.node, "", "  ")
	if err != nil {
		return nil
	}
	return jsonAst
}

// PrintAST 打印当前 meta 的 ast 节点树
func (m *meta) PrintAST() {
	ast.Print(token.NewFileSet(), m.node)
}

// Expression 按照当前 meta 的 ast 节点输出其所在文件中的对应的 代码
func (m *meta) Expression() string {
	fileContent, err := os.ReadFile(m.path)
	if err != nil {
		return ""
	}
	fileContentLen := len(fileContent)
	if m.node.Pos() > m.node.End() || int(m.node.Pos()) >= fileContentLen || int(m.node.End()) > fileContentLen {
		return ""
	}
	return strings.TrimSpace(string(fileContent[m.node.Pos()-1 : m.node.End()-1]))
}

// Format 按照当前 meta 的 ast 节点格式化输出其对应的 代码
func (m *meta) Format() string {
	buffer := &bytes.Buffer{}
	err := format.Node(buffer, token.NewFileSet(), m.node)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
