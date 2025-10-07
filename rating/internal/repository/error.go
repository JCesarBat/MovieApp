package repository

import "errors"

// ErrNotFound is returned when a requested record is not
// found.
var ErrNotFound = errors.New("not Found")

var ErrInvalidRecordType = errors.New("not found this record Type")
