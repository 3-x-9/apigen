package generator

import (
	"fmt"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

type Generator struct{}

type CommandInfo struct {
	GoName  string
	CLIName string
}

func NewGenerator() *Generator {
	return &Generator{}
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (g *Generator) Generate(specPath, outputDir string, moduleName string) error {
	// 1. Load OpenAPI spec
	loader := openapi3.NewLoader()

	var doc *openapi3.T
	var err error

	if isURL(specPath) {
		parsedURL, err := url.Parse(specPath)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}
		doc, err = loader.LoadFromURI(parsedURL)
	} else {
		doc, err = loader.LoadFromFile(specPath)
	}

	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	if err := createFolders(outputDir); err != nil {
		return err
	}

	// 4. Detect Auth
	schemes := detectAuth(doc)

	// 2. Write go.mod
	if err := writeGoMod(outputDir, moduleName); err != nil {
		return err
	}

	// 3. Write config and root cmd
	if err := writeConfig(outputDir, schemes, doc.Servers); err != nil {
		return err
	}
	if err := writeRootCmd(outputDir, moduleName); err != nil {
		return err
	}

	if err := writeModels(outputDir, doc); err != nil {
		return err
	}

	tagToCmds := make(map[string][]CommandInfo)
	for path, pathItem := range doc.Paths.Map() {
		ops := map[string]*openapi3.Operation{
			"get":    pathItem.Get,
			"post":   pathItem.Post,
			"put":    pathItem.Put,
			"delete": pathItem.Delete,
		}

		for method, op := range ops {
			if op == nil {
				continue
			}

			tag := "Misc"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}

			goName := sanitizeCommandName(path, method)
			cliName := sanitizeCLIName(path, method)
			if err := writeEndpointCmd(outputDir, moduleName, goName, cliName, op, path, method, schemes, &doc.Security); err != nil {
				return err
			}

			tagToCmds[tag] = append(tagToCmds[tag], CommandInfo{GoName: goName, CLIName: cliName})
		}
	}

	for tag := range tagToCmds {
		if err := writeTagCmd(outputDir, tag); err != nil {
			return err
		}
	}

	if err := writeMain(outputDir, moduleName, tagToCmds); err != nil {
		return err
	}

	fmt.Printf("Generated CLI in %s\n", outputDir)
	return nil
}
