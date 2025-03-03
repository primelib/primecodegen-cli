package openapipatch

import (
	"errors"
	"fmt"
	"strings"

	"github.com/primelib/primecodegen/pkg/commonpatch"
	"github.com/primelib/primecodegen/pkg/openapi/openapidocument"
	"github.com/primelib/primecodegen/pkg/util"
	"github.com/rs/zerolog/log"
)

func ApplyPatches(input []byte, patches []string) ([]byte, error) {
	for _, patch := range patches {
		// file-based patch
		if strings.HasPrefix(patch, "file:") {
			patchFile := strings.TrimPrefix(patch, "file:")
			patchedBytes, patchErr := commonpatch.ApplyPatchFile(input, patchFile)
			if patchErr != nil {
				return input, errors.Join(util.ErrFailedToPatchDocument, patchErr)
			}

			input = patchedBytes
		}

		// builtin patch
		if p, ok := V3Patchers[patch]; ok {
			log.Debug().Str("id", p.ID).Msg("applying patch to spec")

			doc, err := openapidocument.OpenDocument(input)
			if err != nil {
				return input, err
			}

			v3doc, errs := doc.BuildV3Model()
			if len(errs) > 0 {
				return input, fmt.Errorf("failed to build v3 high level model: %w", errors.Join(errs...))
			}

			patchErr := p.Func(v3doc)
			if patchErr != nil {
				return input, fmt.Errorf("failed to patch document with [%s]: %w", p.ID, patchErr)
			}

			bytes, err := openapidocument.RenderV3Document(v3doc)
			if err != nil {
				return input, errors.Join(util.ErrRenderDocument, err)
			}

			input = bytes
		}
	}

	return input, nil
}
