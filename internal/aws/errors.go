package aws

import "errors"

// Package-level shared errors
var errMissingAccountID = errors.New("aws caller identity response missing account ID")
