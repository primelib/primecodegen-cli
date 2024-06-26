package openapi_go

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
	texttemplate "text/template"

	"github.com/cidverse/go-ptr"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/primelib/primecodegen/pkg/openapi/openapigenerator"
	"github.com/primelib/primecodegen/pkg/template"
	"github.com/primelib/primecodegen/pkg/util"
	"github.com/rs/zerolog/log"
)

type GoGenerator struct {
	reservedWords  []string
	primitiveTypes []string
	typeToImport   map[string]string
}

func (g *GoGenerator) Id() string {
	return "go"
}

func (g *GoGenerator) Description() string {
	return "Generates Go client code"
}

func (g *GoGenerator) Generate(opts openapigenerator.GenerateOpts) error {
	// check opts
	if opts.Doc == nil {
		return fmt.Errorf("document is required")
	}

	// required options
	if opts.ArtifactId == "" {
		return fmt.Errorf("artifact id is required, please set the --md-artifact-id flag")
	}

	// build template data
	templateData, err := g.TemplateData(opts.Doc)
	if err != nil {
		return fmt.Errorf("failed to build template data: %w", err)
	}

	// set packages
	templateData.Packages = openapigenerator.CommonPackages{
		Client:     "client",
		Models:     "models",
		Enums:      "enums",
		Operations: "operations",
		Auth:       "auth",
	}

	// generate files
	files, err := openapigenerator.GenerateFiles(fmt.Sprintf("openapi-%s-%s", g.Id(), opts.TemplateId), opts.OutputDir, templateData, template.RenderOpts{
		DryRun:               opts.DryRun,
		Types:                nil,
		IgnoreFiles:          nil,
		IgnoreFileCategories: nil,
		Properties:           map[string]string{},
		TemplateFunctions: texttemplate.FuncMap{
			"toClassName":     g.ToClassName,
			"toFunctionName":  g.ToFunctionName,
			"toPropertyName":  g.ToPropertyName,
			"toParameterName": g.ToParameterName,
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}
	for _, f := range files {
		log.Debug().Str("file", f.File).Str("template-file", f.TemplateFile).Str("state", string(f.State)).Msg("Generated file")
	}
	log.Info().Msgf("Generated %d files", len(files))

	// post-processing (formatting)
	err = g.PostProcessing(opts.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to run post-processing: %w", err)
	}

	return nil
}

func (g *GoGenerator) TemplateData(doc *libopenapi.DocumentModel[v3.Document]) (openapigenerator.DocumentModel, error) {
	return openapigenerator.BuildTemplateData(doc, g)
}

func (g *GoGenerator) ToClassName(name string) string {
	name = util.ToPascalCase(name)

	if slices.Contains(g.reservedWords, name) {
		return name + "Model"
	}
	return name
}

func (g *GoGenerator) ToFunctionName(name string) string {
	name = util.ToPascalCase(name)

	if slices.Contains(g.reservedWords, name) {
		return name + "Func"
	}

	return name
}

func (g *GoGenerator) ToPropertyName(name string) string {
	name = util.ToPascalCase(name)

	if slices.Contains(g.reservedWords, name) {
		return name + "Prop"
	}

	return name
}

func (g *GoGenerator) ToParameterName(name string) string {
	name = util.ToCamelCase(name)

	if slices.Contains(g.reservedWords, name) {
		return name + "Prop"
	}

	return name
}

func (g *GoGenerator) ToCodeType(schema *base.Schema, required bool) (string, error) {
	// multiple types
	if util.CountExcluding(schema.Type, "null") > 1 {
		return "interface{}", nil
	}

	// normal types
	var codeType string
	switch {
	case slices.Contains(schema.Type, "string"):
		switch schema.Format {
		case "uri", "binary", "byte":
			codeType = "string"
		case "date", "date-time":
			codeType = "string"
		default:
			codeType = "string"
		}
	case slices.Contains(schema.Type, "boolean"):
		codeType = "bool"
	case slices.Contains(schema.Type, "integer"):
		switch schema.Format {
		case "int32":
			codeType = "int32"
		case "int64":
			codeType = "int64"
		default:
			codeType = "int64"
		}
	case slices.Contains(schema.Type, "number"):
		switch schema.Format {
		case "float":
			codeType = "float32"
		case "double":
			codeType = "float64"
		default:
			codeType = "float64"
		}
	case slices.Contains(schema.Type, "array"):
		arrayType, err := g.ToCodeType(schema.Items.A.Schema(), true)
		if err != nil {
			return "", fmt.Errorf("unhandled array type. schema: %s, format: %s, message: %w", schema.Type, schema.Format, err)
		}
		codeType = "[]" + arrayType
	case slices.Contains(schema.Type, "object"):
		if schema.AdditionalProperties != nil {
			keyType := "string"
			additionalProperties, err := g.ToCodeType(schema.AdditionalProperties.A.Schema(), true)
			if err != nil {
				return "", fmt.Errorf("unhandled additional properties type. schema: %s, format: %s: %w", schema.Type, schema.Format, err)
			}
			codeType = "map[" + keyType + "]" + additionalProperties
		} else if schema.AdditionalProperties == nil && schema.Properties == nil {
			codeType = "interface{}"
		} else {
			if schema.Title == "" {
				// TODO: ensure all schemas have a title
				// return "", fmt.Errorf("schema does not have a title. schema: %s", schema.Type)
			}
			codeType = g.ToClassName(schema.Title)
		}
	default:
		return "", fmt.Errorf("unhandled type. schema: %s, format: %s", schema.Type, schema.Format)
	}

	// pointer
	if !required && !strings.HasPrefix(codeType, "[]") && !strings.HasPrefix(codeType, "map[") && ptr.ValueOrDefault(schema.Nullable, true) == true {
		codeType = "*" + codeType
	}

	return codeType, nil
}

func (g *GoGenerator) PostProcessType(codeType string) string {
	return codeType
}

func (g *GoGenerator) IsPrimitiveType(input string) bool {
	return slices.Contains(g.primitiveTypes, input)
}

func (g *GoGenerator) TypeToImport(typeName string) string {
	if typeName == "" {
		return ""
	}
	typeName = strings.Replace(typeName, "*", "", -1)

	return g.typeToImport[typeName]
}

func (g *GoGenerator) PostProcessing(outputDir string) error {
	// run gofmt
	cmd := exec.Command("gofmt", "-s", "-w", outputDir)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running gofmt: %v", err)
	}

	return nil
}

func NewGenerator() *GoGenerator {
	// references: https://openapi-generator.tech/docs/generators/go
	return &GoGenerator{
		reservedWords: []string{
			"bool",
			"break",
			"byte",
			"case",
			"chan",
			"complex128",
			"complex64",
			"const",
			"continue",
			"default",
			"defer",
			"else",
			"error",
			"fallthrough",
			"float32",
			"float64",
			"for",
			"func",
			"go",
			"goto",
			"if",
			"import",
			"int",
			"int16",
			"int32",
			"int64",
			"int8",
			"interface",
			"map",
			"nil",
			"package",
			"range",
			"return",
			"rune",
			"select",
			"string",
			"struct",
			"switch",
			"type",
			"uint",
			"uint16",
			"uint32",
			"uint64",
			"uint8",
			"uintptr",
			"var",
		},
		primitiveTypes: []string{
			"string",
			"bool",
			"int",
			"int32",
			"int64",
			"float32",
			"float64",
			"byte",
			"rune",
			"time.Time",
		},
		typeToImport: map[string]string{
			"time.Time": "time",
		},
	}
}
