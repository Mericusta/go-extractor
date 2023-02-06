package extractor

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
)

type GoImportMeta struct {
	*meta
	alias      string
	importPath string
}

func ExtractGoImportMeta(extractFilepath string, alias string) (*GoImportMeta, error) {
	gfm, err := ExtractGoFileMeta(extractFilepath)
	if err != nil {
		return nil, err
	}

	gim := SearchGoImportMeta(gfm.meta, alias)
	if gim == nil {
		return nil, fmt.Errorf("can not find import node")
	}

	return gim, nil
}

func SearchGoImportMeta(m *meta, alias string) *GoImportMeta {
	var importSpec *ast.ImportSpec
	ast.Inspect(m.node, func(n ast.Node) bool {
		if IsImportNode(n) {
			spec := n.(*ast.ImportSpec)
			if spec.Name != nil && spec.Name.String() == alias {
				importSpec = spec
			} else if filepath.Base(strings.Trim(spec.Path.Value, "\"")) == alias {
				importSpec = spec
			}
		}
		return importSpec == nil
	})
	if importSpec == nil {
		return nil
	}
	return &GoImportMeta{
		meta:       m.newMeta(importSpec),
		alias:      alias,
		importPath: importSpec.Path.Value,
	}
}

func IsImportNode(n ast.Node) bool {
	_, ok := n.(*ast.ImportSpec)
	return ok
}

func (gim *GoImportMeta) Alias() string {
	return gim.alias
}

func (gim *GoImportMeta) ImportPath() string {
	return gim.importPath
}
