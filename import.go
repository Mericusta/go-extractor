package extractor

// type GoImportMeta struct {
// 	*meta
// }

// func ExtractGoImportMeta(extractFilepath, alias string) (*GoImportMeta, error) {
// 	gfm, err := ExtractGoFileMeta(extractFilepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gim := SearchGoImportMeta(gfm.meta, alias)
// 	if gim == nil {
// 		return nil, fmt.Errorf("can not find import node")
// 	}

// 	return gim, nil
// }

// func SearchGoImportMeta(m *meta, alias string) *GoImportMeta {
// 	var importSpec *ast.ImportSpec
// 	ast.Inspect(m.node, func(n ast.Node) bool {
// 		if IsImportNode(n) {
// 			spec := n.(*ast.ImportSpec)
// 			if spec.Name != nil && spec.Name.String() == alias {
// 				importSpec = spec
// 			} else if filepath.Base(strings.Trim(spec.Path.Value, "\"")) == alias {
// 				importSpec = spec
// 			}
// 		}
// 		return importSpec == nil
// 	})
// 	if importSpec == nil {
// 		return nil
// 	}
// 	return &GoImportMeta{
// 		meta: m.copyMeta(importSpec),
// 	}
// }

// func (gim *GoImportMeta) Alias() string {
// 	aliasIdent := gim.node.(*ast.ImportSpec).Name
// 	if aliasIdent == nil {
// 		return ""
// 	}
// 	return aliasIdent.String()
// }

// func (gim *GoImportMeta) ImportPath() string {
// 	return filepath.Base(strings.Trim(gim.node.(*ast.ImportSpec).Path.Value, "\""))
// }
