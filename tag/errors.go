package tag

import "github.com/mangaweb4/mangaweb4-backend/errors"

var ErrTagNotFound = errors.New(1_000_000, "tag '%s' not found.")
