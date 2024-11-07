package openapipatch

import (
	"testing"

	"github.com/primelib/primecodegen/pkg/openapi/openapidocument"
	"github.com/stretchr/testify/assert"
)

func TestPruneOperationTags(t *testing.T) {
	// parse spec
	const spec = `
    openapi: 3.0.0
    info:
      title: Sample API
      version: 1.0.0
    paths:
      /test:
        get:
          tags:
            - test
          responses:
            '200':
              description: OK
    `
	document, err := openapidocument.OpenDocument([]byte(spec))
	if err != nil {
		t.Fatalf("error creating document: %v", err)
	}
	v3doc, errors := document.BuildV3Model()
	assert.Equal(t, 0, len(errors))

	// prune operation tags
	_ = PruneOperationTags(v3doc)

	// check if tags are pruned
	v, _ := v3doc.Model.Paths.PathItems.Get("/test")
	assert.Nil(t, v.Get.Tags, "tags should be pruned")
}

func TestInlineAllOfHierarchies(t *testing.T) {
	// arrange
	const spec = `
    openapi: 3.0.0
    info:
      title: Sample API
      version: v1.0.0
    components:
        schemas:
            TestSchema:
                properties:
                  propertyA:
                    type: string
                    description: Description A
                  propertyB:
                    type: string
                    description: Description B
            TestSchemaWithReferences:
                allOf:
                    - $ref: '#/components/schemas/TestSchema'
                    - properties:
                        additionalPropertyC:
                          type: string
                          description: Description C
    `
	document, err := openapidocument.OpenDocument([]byte(spec))
	if err != nil {
		t.Fatalf("error creating document: %v", err)
	}
	v3doc, errors := document.BuildV3Model()
	assert.Equal(t, 0, len(errors))

	// act
	_ = InlineAllOfHierarchies(v3doc)

	// assert
	_, document, _, errors = document.RenderAndReload()
	assert.Equal(t, 0, len(errors))
	v3doc, errors = document.BuildV3Model()
	assert.Equal(t, 0, len(errors))

	propsToCheck := []string{"propertyA", "propertyB", "additionalPropertyC"}

	for schemaMapEntry := v3doc.Model.Components.Schemas.Oldest(); schemaMapEntry != nil; schemaMapEntry = schemaMapEntry.Next() {
		schema, err := schemaMapEntry.Value.BuildSchema()
		assert.NoError(t, err)
		assert.Nil(t, schema.AllOf, "allOf references should be deleted")
		if schemaMapEntry.Key == "TestSchemaWithReferences" {
			assert.Equal(t, 3, schema.Properties.Len())
			for _, prop := range propsToCheck {
				_, exists := schema.Properties.Get(prop)
				assert.True(t, exists, "Property \"%s\" is missing!", prop)
			}
		}
	}
}
