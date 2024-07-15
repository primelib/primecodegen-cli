package openapigenerator

import (
	"fmt"
	"os"
	"strings"

	"github.com/primelib/primecodegen/pkg/template"
	"github.com/primelib/primecodegen/pkg/util"
	"gopkg.in/yaml.v3"
)

func GeneratorById(id string, allGenerators []CodeGenerator) (CodeGenerator, error) {
	for _, g := range allGenerators {
		if g.Id() == id {
			return g, nil
		}
	}

	return nil, fmt.Errorf("generator with id %s not found", id)
}

func GenerateFiles(templateId string, outputDir string, templateData DocumentModel, renderOpts template.RenderOpts, generatorOpts GenerateOpts) ([]template.RenderedFile, error) {
	var files []template.RenderedFile

	// print template data
	if os.Getenv("PRIMECODEGEN_DEBUG_TEMPLATEDATA") == "true" {
		bytes, _ := yaml.Marshal(templateData)
		fmt.Print(string(bytes))
	}

	// global template data
	common := GlobalTemplate{
		GeneratorProperties: renderOpts.Properties,
		Endpoints:           templateData.Endpoints,
		Auth:                templateData.Auth,
		Packages:            templateData.Packages,
		Services:            templateData.Services,
		Operations:          templateData.Operations,
		Models:              templateData.Models,
		Enums:               templateData.Enums,
	}
	metadata := Metadata{
		ArtifactGroupId: generatorOpts.ArtifactGroupId,
		ArtifactId:      generatorOpts.ArtifactId,
		Name:            strings.TrimSpace(templateData.Name),
		DisplayName:     strings.TrimSpace(templateData.DisplayName),
		Description:     templateData.Description,
		RepositoryUrl:   generatorOpts.RepositoryUrl,
		LicenseName:     generatorOpts.LicenseName,
		LicenseUrl:      generatorOpts.LicenseUrl,
	}
	if metadata.ArtifactId == "" {
		metadata.ArtifactId = util.ToSlug(metadata.Name)
	}

	var data []interface{}
	data = append(data, SupportOnceTemplate{
		Metadata: metadata,
		Common:   common,
	})
	data = append(data, APIOnceTemplate{
		Metadata: metadata,
		Common:   common,
		Package:  common.Packages.Client,
	})

	for _, service := range templateData.Services {
		data = append(data, APIEachTemplate{
			Metadata:       metadata,
			Common:         common,
			Package:        common.Packages.Client,
			TagName:        service.Name,
			TagDescription: service.Description,
			TagOperations:  service.Operations,
		})
	}
	for _, op := range templateData.Operations {
		data = append(data, OperationEachTemplate{
			Metadata:  metadata,
			Common:    common,
			Package:   common.Packages.Operations,
			Name:      op.Name,
			Operation: op,
		})
	}
	for _, model := range templateData.Models {
		data = append(data, ModelEachTemplate{
			Metadata: metadata,
			Common:   common,
			Package:  common.Packages.Models,
			Name:     model.Name,
			Model:    model,
		})
	}
	for _, enum := range templateData.Enums {
		data = append(data, EnumEachTemplate{
			Metadata: metadata,
			Common:   common,
			Package:  common.Packages.Models,
			Name:     enum.Name,
			Enum:     enum,
		})
	}

	// render files
	for _, d := range data {
		var renderedFiles []template.RenderedFile
		var renderErr error

		if _, ok := d.(SupportOnceTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeSupportOnce, d, renderOpts)
		}
		if _, ok := d.(APIOnceTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeAPIOnce, d, renderOpts)
		}
		if _, ok := d.(APIEachTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeAPIEach, d, renderOpts)
		}
		if _, ok := d.(OperationEachTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeOperationEach, d, renderOpts)
		}
		if _, ok := d.(ModelEachTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeModelEach, d, renderOpts)
		}
		if _, ok := d.(EnumEachTemplate); ok {
			renderedFiles, renderErr = template.RenderTemplateById(templateId, outputDir, template.TypeEnumEach, d, renderOpts)
		}

		if renderErr != nil {
			return nil, fmt.Errorf("failed to render template: %w", renderErr)
		}
		files = append(files, renderedFiles...)
	}

	return files, nil
}
