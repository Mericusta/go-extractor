package extractor

import (
	"os"
)

func ReadUnitTestFile(p string) []byte {
	c, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return c
}
