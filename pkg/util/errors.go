package util

import (
	"fmt"
)

var (
	ErrFileMissing           = fmt.Errorf("file is missing")
	ErrNoFilesSpecified      = fmt.Errorf("no files specified")
	ErrDocumentMerge         = fmt.Errorf("failed to merge documents")
	ErrFailedToPatchDocument = fmt.Errorf("failed to patch document")
	ErrWriteDocumentToFile   = fmt.Errorf("failed to write document to file")
	ErrWriteDocumentToStdout = fmt.Errorf("failed to write document to stdout")
	ErrJSONMarshal           = fmt.Errorf("failed to marshal into JSON")
)
