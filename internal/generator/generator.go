package generator

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(specPath, outputDir string) error {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}
	paths := doc.Paths
	fmt.Printf("Loaded %d paths \n", paths.Len())
	return nil
}
