package meta

import "github.com/mangaweb4/mangaweb4-backend/errors"

var ErrMetaDataNotFound = errors.New(2_000_000, "metadata for '%s' not found.")
