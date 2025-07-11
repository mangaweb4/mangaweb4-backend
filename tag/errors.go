package tag

import "github.com/mangaweb4/mangaweb4-backend/errors"

var ErrTagNotFound = errors.New(1_000_000, "tag '%s' not found.")
var ErrInvalidTagFilter = errors.New(1_000_001, "filter '%s' is invalid.")
var ErrInvalidTagSortField = errors.New(1_000_002, "sort field '%s' is invalid.")
