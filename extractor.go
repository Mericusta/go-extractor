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

func (m *meta) Copy() *meta {
	// var buf bytes.Buffer
	// if err := gob.NewEncoder(&buf).Encode(m); err != nil {
	// 	fmt.Printf("gob.NewEncoder.Encode meta occurs error: %v\n", err)
	// 	return nil
	// }
	dst := &meta{}
	// if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst); err != nil {
	// 	fmt.Printf("gob.NewDecoder.Decode meta occurs error: %v\n", err)
	// 	return nil
	// }
	return dst

	// var root ast.Node

	// ast.Inspect(m.node, func(n ast.Node) bool {
	// 	switch n.(type) {

	// 	}
	// })
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
	if m.node.Pos() > m.node.End() || int(m.node.Pos()) >= fileContentLen || int(m.node.End()) > fileContentLen {
		return ""
	}
	return strings.TrimSpace(string(fileContent[m.node.Pos()-1 : m.node.End()-1]))
}
