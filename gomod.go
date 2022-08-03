package extractor

import (
	"fmt"
	"os"
	"regexp"
)

var (
	GO_MODULE_EXPRESSION_SUB_EXPRESSION_NAME = `NAME`
	GO_MODULE_EXPRESSION                     = fmt.Sprintf(`module\s+(?P<%v>\S+)`, GO_MODULE_EXPRESSION_SUB_EXPRESSION_NAME)
)

func extractGoModuleName(goModFilePath string) (string, error) {
	goModFileContent, err := os.ReadFile(goModFilePath)
	if err != nil {
		return "", err
	}

	goModuleRegexp := regexp.MustCompile(GO_MODULE_EXPRESSION)
	if goModuleRegexp == nil {
		return "", fmt.Errorf("go module regexp is nil")
	}

	nameIndex := goModuleRegexp.SubexpIndex(GO_MODULE_EXPRESSION_SUB_EXPRESSION_NAME)
	submatchSlice := goModuleRegexp.FindSubmatch(goModFileContent)
	if nameIndex >= len(submatchSlice) {
		return "", fmt.Errorf("can not find go module submatch name")
	}
	return string(submatchSlice[nameIndex]), nil
}
