package query

import "github.com/ginger-core/errors"

var (
	ModelNotDefinedErr = errors.Internal().
		WithDetail(errors.Detail{
			"error": "model is not defined to handle query for",
		})
)
