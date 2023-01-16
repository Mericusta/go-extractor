package extractor

import (
	"go/ast"
	"go/token"
	"os"
	"strings"
)

type Meta interface {
	PrintAST()
	Expression() string
}

type meta struct {
	node ast.Node // ast node
	name string   // node from file name
	path string   // node from file path
}

func (m *meta) newMeta(node ast.Node) *meta {
	return &meta{node: node, name: m.name, path: m.path}
}

func (m *meta) PrintAST() {
	ast.Print(token.NewFileSet(), m.node)
}

func (m *meta) Expression() string {
	fileContent, err := os.ReadFile(m.path)
	if err != nil {
		return ""
	}
	fileContentLen := len(fileContent)
	if m.node.Pos() > m.node.End() || int(m.node.Pos()) >= fileContentLen || int(m.node.End()) >= fileContentLen {
		return ""
	}
	return strings.TrimSpace(string(fileContent[m.node.Pos()-1 : m.node.End()-1]))
}
